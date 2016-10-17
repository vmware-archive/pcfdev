package main_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("PCF Dev provision", func() {
	var (
		dockerID string
		pwd      string
	)

	BeforeEach(func() {
		var err error
		pwd, err = os.Getwd()
		Expect(err).NotTo(HaveOccurred())

		output, err := exec.Command("docker", "run", "--privileged", "-d", "-w", "/go/src/pcfdev", "-v", pwd+":/go/src/pcfdev", "pcfdev/provision", "bash", "-c", "umount /etc/resolv.conf && sleep 1000").Output()
		Expect(err).NotTo(HaveOccurred())
		dockerID = strings.TrimSpace(string(output))

		Expect(exec.Command("bash", "-c", "echo \"#!/bin/bash\necho 'Waiting for services to start...'\necho \\$@\" > "+pwd+"/provision-script").Run()).To(Succeed())
		Expect(exec.Command("docker", "exec", dockerID, "chmod", "+x", "/go/src/pcfdev/provision-script").Run()).To(Succeed())
		Expect(exec.Command("docker", "exec", dockerID, "go", "build", "-ldflags", "-X main.provisionScriptPath=/go/src/pcfdev/provision-script", "pcfdev").Run()).To(Succeed())
	})

	AfterEach(func() {
		os.RemoveAll(filepath.Join(pwd, "pcfdev"))
		os.RemoveAll(filepath.Join(pwd, "provision-script"))
		Expect(exec.Command("docker", "rm", dockerID, "-f").Run()).To(Succeed())
	})

	It("should provision PCF Dev", func() {
		session, err := gexec.Start(exec.Command("docker", "exec", dockerID, "/go/src/pcfdev/pcfdev", "local.pcfdev.io", "192.168.11.11"), GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session, "10s").Should(gexec.Exit(0))
		Expect(session).To(gbytes.Say("Waiting for services to start..."))
	})

	It("should create certificates", func() {
		session, err := gexec.Start(exec.Command("docker", "exec", dockerID, "/go/src/pcfdev/pcfdev", "local.pcfdev.io", "192.168.11.11"), GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session, "10s").Should(gexec.Exit(0))

		session, err = gexec.Start(exec.Command("docker", "exec", dockerID, "bash", "-c", "echo 127.0.0.1 local.pcfdev.io >> /etc/hosts"), GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))

		session, err = gexec.Start(exec.Command("docker", "exec", "-d", dockerID, "service", "nginx", "start"), GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))

		session, err = gexec.Start(exec.Command("docker", "exec", dockerID, "curl", "--cacert", "/var/pcfdev/openssl/ca_cert.pem", "https://local.pcfdev.io:443"), GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))
	})

	It("should disable HSTS in UAA", func() {
		session, err := gexec.Start(exec.Command("docker", "exec", dockerID, "/go/src/pcfdev/pcfdev", "local.pcfdev.io", "192.168.11.11"), GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session, "10s").Should(gexec.Exit(0))

		session, err = gexec.Start(exec.Command("docker", "exec", dockerID, "grep", "-A", "1", "<param-name>hstsEnabled</param-name>", "/var/vcap/packages/uaa/tomcat/conf/web.xml"), GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))
		Eventually(session).Should(gbytes.Say("<param-value>false</param-value>"))
	})

	It("should resolve *.cf.internal to the internal IP using Dnsmasq", func() {
		session, err := gexec.Start(exec.Command("docker", "exec", dockerID, "/go/src/pcfdev/pcfdev", "local.pcfdev.io", "192.168.11.11"), GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session, "10s").Should(gexec.Exit(0))

		session, err = gexec.Start(exec.Command("docker", "exec", dockerID, "host", "bbs.service.cf.internal"), GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))
		Eventually(session).Should(gbytes.Say(`bbs.service.cf.internal has address [0-9]+\.[0-9]+\.[0-9]+\.[0-9]`))
	})

	Context("when the distribution is not 'pcf'", func() {
		BeforeEach(func() {
			Expect(exec.Command("docker", "exec", dockerID, "go", "build", "-ldflags", "-X main.distro=oss -X main.provisionScriptPath=/go/src/pcfdev/provision-script", "pcfdev").Run()).To(Succeed())
		})

		It("should not disable HSTS in UAA", func() {
			session, err := gexec.Start(exec.Command("docker", "exec", dockerID, "/go/src/pcfdev/pcfdev", "local.pcfdev.io", "192.168.11.11"), GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session, "10s").Should(gexec.Exit(0))

			session, err = gexec.Start(exec.Command("docker", "exec", dockerID, "grep", "<param-name>hstsEnabled</param-name>", "/var/vcap/packages/uaa/tomcat/conf/web.xml"), GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit(1))
		})
	})

	Context("when provisioning fails", func() {
		BeforeEach(func() {
			Expect(exec.Command("bash", "-c", "echo \"#!/bin/bash\nexit 42\" > "+pwd+"/provision-script").Run()).To(Succeed())
		})

		It("should exit with the exit status of the provision script", func() {
			session, _ := gexec.Start(exec.Command("docker", "exec", dockerID, "/go/src/pcfdev/pcfdev", "local.pcfdev.io", "192.168.11.11"), GinkgoWriter, GinkgoWriter)
			Eventually(session, "10s").Should(gexec.Exit(42))
		})
	})

	Context("when provisioning takes too long", func() {
		BeforeEach(func() {
			Expect(exec.Command("bash", "-c", "echo \"#!/bin/bash\nsleep 20\" > "+pwd+"/provision-script").Run()).To(Succeed())
			Expect(exec.Command("docker", "exec", dockerID, "go", "build", "-ldflags", "-X main.provisionScriptPath=/go/src/pcfdev/provision-script -X main.timeoutInSeconds=2", "pcfdev").Run()).To(Succeed())
		})

		It("exit with an exit status of 1 and tell why it is exiting...", func() {
			session, err := gexec.Start(exec.Command("docker", "exec", dockerID, "/go/src/pcfdev/pcfdev", "local.pcfdev.io", "192.168.11.11"), GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session, "15s").Should(gexec.Exit(1))
			Expect(session).To(gbytes.Say("Timed out after 2 seconds."))
		})
	})
})
