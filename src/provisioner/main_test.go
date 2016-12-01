package main_test

import (
	"os"
	"os/exec"
	"strings"

	"regexp"

	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"net"
	"time"
	"strconv"
	"math/rand"
)

var _ = Describe("PCF Dev provision", func() {
	var (
		dockerID string
		pwd      string
		dockerExecTimeout string = "3s"
		randomTcpPort string
	)

	var portsForwarded struct {
		Port80        string
		Port8060      string
		Port443       string
		Port22        string
		Port2222      string
		Port61001     string
		Port61100     string
		RandomTcpPort string

	}


	BeforeEach(func() {
		portsForwarded.Port80 = randomOpenPort()
		portsForwarded.Port8060 = randomOpenPort()
		portsForwarded.Port443 = randomOpenPort()
		portsForwarded.Port22 = randomOpenPort()
		portsForwarded.Port2222 = randomOpenPort()
		portsForwarded.Port61001 = randomOpenPort()
		portsForwarded.Port61100= randomOpenPort()
		portsForwarded.RandomTcpPort = randomOpenPort()

		var err error
		pwd, err = os.Getwd()
		Expect(err).NotTo(HaveOccurred())

		randomTcpPort = randomPortInRange("61001", "61100")

		output, err := exec.Command(
			"docker", "run",
			"-p", portsForwarded.Port80+":80",
			"-p", portsForwarded.Port8060 +":8060",
			"-p", portsForwarded.Port443+":443",
			"-p", portsForwarded.Port22+":22",
			"-p", portsForwarded.Port2222+":2222",
			"-p", portsForwarded.Port61001+":61001",
			"-p", portsForwarded.Port61100+":61100",
			"-p", portsForwarded.RandomTcpPort +":"+ randomTcpPort,
			"--privileged",
			"-d",
			"-v", pwd+":/go/src/provisioner",
			"-w", "/go/src/provisioner",
			"pcfdev/provision",
			"bash", "-c", "umount /etc/resolv.conf && sleep 1000",
		).Output()
		Expect(err).NotTo(HaveOccurred())
		dockerID = strings.TrimSpace(string(output))

		Expect(exec.Command("bash", "-c", "echo \"#!/bin/bash\necho 'Waiting for services to start...'\necho \\$@\" > "+pwd+"/provision-script").Run()).To(Succeed())
		Expect(exec.Command("docker", "exec", dockerID, "chmod", "+x", "/go/src/provisioner/provision-script").Run()).To(Succeed())
		Expect(exec.Command("docker", "exec", dockerID, "go", "build", "-ldflags", "-X main.provisionScriptPath=/go/src/provisioner/provision-script", "-o", "provision", "provisioner").Run()).To(Succeed())
	})

	AfterEach(func() {
		os.RemoveAll(filepath.Join(pwd, "provision"))
		os.RemoveAll(filepath.Join(pwd, "provision-script"))
		Expect(exec.Command("docker", "rm", dockerID, "-f").Run()).To(Succeed())
	})

	It("should provision PCF Dev", func() {
		session := provisionForVirtualBox(dockerID)
		Expect(session).To(gbytes.Say("Waiting for services to start..."))

		session = runSuccessfully(exec.Command("docker", "exec", dockerID, "file", "/run/pcfdev-healthcheck"), dockerExecTimeout)
		Expect(session).NotTo(gbytes.Say("No such file or directory"))
	})

	It("should set up the monitrc files for an HTTP server and an root executable _ctl script running on the box", func() {
		provisionForVirtualBox(dockerID)

		session := runSuccessfully(exec.Command("docker", "exec", dockerID, "cat", "/var/vcap/monit/job/1001_pcfdev_api.monitrc"), dockerExecTimeout)
		Expect(session).To(gbytes.Say("check process pcfdev-api"))
		Expect(session).To(gbytes.Say("with pidfile /var/vcap/sys/run/pcfdev-api/api.pid"))
		Expect(session).To(gbytes.Say(`start program "/var/pcfdev/api/api_ctl start"`))
		Expect(session).To(gbytes.Say(`stop program "/var/pcfdev/api/api_ctl stop"`))
		Expect(session).To(gbytes.Say("group vcap"))
		Expect(session).To(gbytes.Say("mode manual"))

		session = runSuccessfully(exec.Command("docker", "exec", dockerID, "ls", "-l", "/var/pcfdev/api/api_ctl"), dockerExecTimeout)

		Expect(session).To(gbytes.Say("-rwxr--r--"))
	})

	It("should create certificates", func() {
		provisionForVirtualBox(dockerID)
		runSuccessfully(exec.Command("docker", "exec", dockerID, "bash", "-c", "echo 127.0.0.1 local.pcfdev.io >> /etc/hosts"), dockerExecTimeout)
		runSuccessfully(exec.Command("docker", "exec", "-d", dockerID, "service", "nginx", "start"), dockerExecTimeout)
		runSuccessfully(exec.Command("docker", "exec", dockerID, "curl", "--cacert", "/var/pcfdev/openssl/ca_cert.pem", "https://local.pcfdev.io:443"), dockerExecTimeout)
	})

	It("should disable HSTS in UAA", func() {
		provisionForVirtualBox(dockerID)

		session := runSuccessfully(exec.Command("docker", "exec", dockerID, "grep", "-A", "1", "<param-name>hstsEnabled</param-name>", "/var/vcap/packages/uaa/tomcat/conf/web.xml"), dockerExecTimeout)
		Eventually(session).Should(gbytes.Say("<param-value>false</param-value>"))
	})

	It("should directly insert the internal-ip into the dns_server flag of garden", func() {
		provisionForVirtualBox(dockerID)

		output, err := exec.Command("docker", "exec", dockerID, "ip", "route", "get", "1").Output()
		Expect(err).NotTo(HaveOccurred())
		regex := regexp.MustCompile(`\s{2}src\s(.*)`)
		internalIP := regex.FindStringSubmatch(string(output))[1]

		session := runSuccessfully(exec.Command("docker", "exec", dockerID, "grep", internalIP, "/var/vcap/jobs/garden/bin/garden_ctl"), dockerExecTimeout)
		Eventually(session).Should(gbytes.Say("-dnsServer=" + internalIP))
	})

	Describe("Network access", func() {
		Context("on AWS", func() {
			It("does not allow connections by default", func() {
				Expect(exec.Command("docker", "exec", "-d", dockerID, "go", "run", "assets/stub_server.go", "8060").Run()).To(Succeed())
				waitForServer("localhost:"+portsForwarded.Port8060, 5*time.Second)
				provisionForAws(dockerID)
				runFailure(exec.Command("curl", "localhost:"+portsForwarded.Port8060), "5s")
			})

			It("allows external access to http cf router", func() {
				Expect(exec.Command("docker", "exec", "-d", dockerID, "go", "run", "assets/stub_server.go", "80").Run()).To(Succeed())
				waitForServer("localhost:"+portsForwarded.Port80, 5*time.Second)
				provisionForAws(dockerID)

				session := runSuccessfully(exec.Command("curl", "localhost:"+portsForwarded.Port80), "5s")
				Expect(session).To(gbytes.Say("Response from port 80 stub server"))
			})

			It("allows external access to https cf router", func() {
				Expect(exec.Command("docker", "exec", "-d", dockerID, "go", "run", "assets/stub_server.go", "443").Run()).To(Succeed())
				waitForServer("localhost:"+portsForwarded.Port443, 5*time.Second)
				provisionForAws(dockerID)

				session := runSuccessfully(exec.Command("curl", "localhost:"+portsForwarded.Port443), "5s")
				Expect(session).To(gbytes.Say("Response from port 443 stub server"))
			})

			It("allows external access to ssh", func() {
				Expect(exec.Command("docker", "exec", "-d", dockerID, "go", "run", "assets/stub_server.go", "22").Run()).To(Succeed())
				waitForServer("localhost:"+portsForwarded.Port22, 5*time.Second)
				provisionForAws(dockerID)

				session := runSuccessfully(exec.Command("curl", "localhost:"+portsForwarded.Port22), "5s")
				Expect(session).To(gbytes.Say("Response from port 22 stub server"))
			})

			It("allows external access to the ssh proxy", func() {
				Expect(exec.Command("docker", "exec", "-d", dockerID, "go", "run", "assets/stub_server.go", "2222").Run()).To(Succeed())
				waitForServer("localhost:"+portsForwarded.Port2222, 5*time.Second)
				provisionForAws(dockerID)
				session := runSuccessfully(exec.Command("curl", "localhost:"+portsForwarded.Port2222), "5s")
				Expect(session).To(gbytes.Say("Response from port 2222 stub server"))
			})

			It("allow external access to tcp router port for lowest port in range", func() {
				Expect(exec.Command("docker", "exec", "-d", dockerID, "go", "run", "assets/stub_server.go", "61001").Run()).To(Succeed())
				waitForServer("localhost:"+portsForwarded.Port61001, 5*time.Second)
				provisionForAws(dockerID)
				session := runSuccessfully(exec.Command("curl", "localhost:"+portsForwarded.Port61001), "5s")
				Expect(session).To(gbytes.Say("Response from port 61001 stub server"))

			})

			It("allow external access to tcp router port for highest port in range", func() {
				Expect(exec.Command("docker", "exec", "-d", dockerID, "go", "run", "assets/stub_server.go", "61100").Run()).To(Succeed())
				waitForServer("localhost:"+portsForwarded.Port61100, 5*time.Second)
				provisionForAws(dockerID)
				session := runSuccessfully(exec.Command("curl", "localhost:"+portsForwarded.Port61100), "5s")
				Expect(session).To(gbytes.Say("Response from port 61100 stub server"))

			})

			It("allow external access to tcp router port for random port in range", func() {
				Expect(exec.Command("docker", "exec", "-d", dockerID, "go", "run", "assets/stub_server.go", randomTcpPort).Run()).To(Succeed())
				waitForServer("localhost:"+portsForwarded.RandomTcpPort, 5*time.Second)
				provisionForAws(dockerID)
				session := runSuccessfully(exec.Command("curl", "localhost:"+portsForwarded.RandomTcpPort), "5s")
				Expect(session).To(gbytes.Say("Response from port "+ randomTcpPort + " stub server"))
			})
		})

		Context("on Virtualbox", func() {
			It("allows connections by default", func() {
				Expect(exec.Command("docker", "exec", "-d", dockerID, "go", "run", "assets/stub_server.go", "8060").Run()).To(Succeed())
				waitForServer("localhost:"+portsForwarded.Port8060, 5*time.Second)
				Expect(exec.Command("docker", "exec", "-d", dockerID, "go", "run", "assets/stub_server.go", "80").Run()).To(Succeed())
				waitForServer("localhost:"+portsForwarded.Port80, 5*time.Second)
				Expect(exec.Command("docker", "exec", "-d", dockerID, "go", "run", "assets/stub_server.go", "443").Run()).To(Succeed())
				waitForServer("localhost:"+portsForwarded.Port443, 5*time.Second)
				Expect(exec.Command("docker", "exec", "-d", dockerID, "go", "run", "assets/stub_server.go", "22").Run()).To(Succeed())
				waitForServer("localhost:"+portsForwarded.Port22, 5*time.Second)
				Expect(exec.Command("docker", "exec", "-d", dockerID, "go", "run", "assets/stub_server.go", "2222").Run()).To(Succeed())
				waitForServer("localhost:"+portsForwarded.Port2222, 5*time.Second)
				Expect(exec.Command("docker", "exec", "-d", dockerID, "go", "run", "assets/stub_server.go", "61001").Run()).To(Succeed())
				waitForServer("localhost:"+portsForwarded.Port61001, 5*time.Second)
				Expect(exec.Command("docker", "exec", "-d", dockerID, "go", "run", "assets/stub_server.go", "61100").Run()).To(Succeed())
				waitForServer("localhost:"+portsForwarded.Port61100, 5*time.Second)
				Expect(exec.Command("docker", "exec", "-d", dockerID, "go", "run", "assets/stub_server.go", randomTcpPort).Run()).To(Succeed())
				waitForServer("localhost:"+portsForwarded.RandomTcpPort, 5*time.Second)

				provisionForVirtualBox(dockerID)

				Expect(runSuccessfully(exec.Command("curl", "localhost:"+portsForwarded.Port80), "5s")).To(gbytes.Say("Response from port 80 stub server"))
				Expect(runSuccessfully(exec.Command("curl", "localhost:"+portsForwarded.Port8060), "5s")).To(gbytes.Say("Response from port 8060 stub server"))
				Expect(runSuccessfully(exec.Command("curl", "localhost:"+portsForwarded.Port443), "5s")).To(gbytes.Say("Response from port 443 stub server"))
				Expect(runSuccessfully(exec.Command("curl", "localhost:"+portsForwarded.Port22), "5s")).To(gbytes.Say("Response from port 22 stub server"))
				Expect(runSuccessfully(exec.Command("curl", "localhost:"+portsForwarded.Port2222), "5s")).To(gbytes.Say("Response from port 2222 stub server"))
				Expect(runSuccessfully(exec.Command("curl", "localhost:"+portsForwarded.Port61001), "5s")).To(gbytes.Say("Response from port 61001 stub server"))
				Expect(runSuccessfully(exec.Command("curl", "localhost:"+portsForwarded.Port61100), "5s")).To(gbytes.Say("Response from port 61100 stub server"))
				Expect(runSuccessfully(exec.Command("curl", "localhost:"+portsForwarded.RandomTcpPort), "5s")).To(gbytes.Say("Response from port "+ randomTcpPort+ " stub server"))
			})
		})
	})

	It("should resolve *.cf.internal to localhost using Dnsmasq", func() {
		provisionForVirtualBox(dockerID)
		session := runSuccessfully(exec.Command("docker", "exec", dockerID, "host", "bbs.service.cf.internal"), "3s")
		Eventually(session).Should(gbytes.Say(`bbs.service.cf.internal has address 127.0.0.1`))
	})

	Context("when the distribution is not 'pcf'", func() {
		BeforeEach(func() {
			Expect(exec.Command("docker", "exec", dockerID, "go", "build", "-ldflags", "-X main.distro=oss -X main.provisionScriptPath=/go/src/provisioner/provision-script", "-o", "provision", "provisioner").Run()).To(Succeed())
		})

		It("should not disable HSTS in UAA", func() {
			provisionForVirtualBox(dockerID)

			session, err := gexec.Start(exec.Command("docker", "exec", dockerID, "grep", "<param-name>hstsEnabled</param-name>", "/var/vcap/packages/uaa/tomcat/conf/web.xml"), GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit(1))
		})
	})

	Context("when provisioning does not have the required number of args", func() {
		It("should exit with an error", func() {
			runFailure(exec.Command("docker", "exec", dockerID, "/go/src/provisioner/provision", "local.pcfdev.io", "192.168.11.11", "", ""), "10s")
		})
	})

	Context("when provisioning fails", func() {
		BeforeEach(func() {
			Expect(exec.Command("bash", "-c", "echo \"#!/bin/bash\nexit 42\" > "+pwd+"/provision-script").Run()).To(Succeed())
		})

		It("should exit with the exit status of the provision script", func() {
			session, _ := gexec.Start(exec.Command("docker", "exec", dockerID, "/go/src/provisioner/provision", "local.pcfdev.io", "192.168.11.11", "", "", "virtualbox"), GinkgoWriter, GinkgoWriter)
			Eventually(session, "10s").Should(gexec.Exit(42))
		})
	})

	Context("when provisioning takes too long", func() {
		BeforeEach(func() {
			Expect(exec.Command("bash", "-c", "echo \"#!/bin/bash\nsleep 20\" > "+pwd+"/provision-script").Run()).To(Succeed())
			Expect(exec.Command("docker", "exec", dockerID, "go", "build", "-ldflags", "-X main.provisionScriptPath=/go/src/provisioner/provision-script -X main.timeoutInSeconds=2", "-o", "provision", "provisioner").Run()).To(Succeed())
		})

		It("exit with an exit status of 1 and tell why it is exiting...", func() {
			session, err := gexec.Start(exec.Command("docker", "exec", dockerID, "/go/src/provisioner/provision", "local.pcfdev.io", "192.168.11.11", "", "", "virtualbox"), GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session, "15s").Should(gexec.Exit(1))
			Expect(session).To(gbytes.Say("Timed out after 2 seconds."))
		})
	})
})

func provisionForVirtualBox(dockerID string) *gexec.Session {
	return runSuccessfully(exec.Command("docker", "exec", dockerID, "/go/src/provisioner/provision", "local.pcfdev.io", "192.168.11.11", "", "", "virtualbox"), "10s")
}

func provisionForAws(dockerID string) *gexec.Session {
	return runSuccessfully(exec.Command("docker", "exec", dockerID, "/go/src/provisioner/provision", "local.pcfdev.io", "192.168.11.11", "", "", "aws"), "10s")
}

func runSuccessfully(command *exec.Cmd, timeout string) *gexec.Session {
	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	EventuallyWithOffset(1, session, timeout).Should(gexec.Exit(0))
	return session
}

func runFailure(command *exec.Cmd, timeout string) *gexec.Session {
	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	ConsistentlyWithOffset(1, session, timeout).ShouldNot(gexec.Exit(0))
	return session
}

func randomOpenPort() string {
	conn, err := net.Listen("tcp", "127.0.0.1:0")
	Expect(err).NotTo(HaveOccurred())
	defer conn.Close()
	address := strings.Split(conn.Addr().String(), ":")
	return address[1]
}

func waitForServer(host string, timeout time.Duration) {
	currentWait := 0 * time.Second
	serverOpen := false

	for !serverOpen && currentWait < timeout {
		exec.Command("curl", host)
		session, _ := gexec.Start(exec.Command("curl", host), GinkgoWriter, GinkgoWriter)
		Eventually(session, "1s").Should(gexec.Exit())
		if session.ExitCode() == 0 {
			serverOpen = true
		}
		currentWait += time.Second * 1
		time.Sleep(time.Second)
	}
	Expect(serverOpen).To(BeTrue())
}

func randomPortInRange(lowerPort string, higherPort string) string {

	higherPortNumber, err := strconv.Atoi(higherPort)
	Expect(err).NotTo(HaveOccurred())

	lowerPortNumber, err := strconv.Atoi(lowerPort)
	Expect(err).NotTo(HaveOccurred())

	return strconv.Itoa(rand.Intn(higherPortNumber - lowerPortNumber) + lowerPortNumber)

}