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

var (
	dockerID   string
	pwd        string
	binaryPath string
)

var _ = BeforeSuite(func() {
	var err error
	pwd, err = os.Getwd()
	Expect(err).NotTo(HaveOccurred())

	output, err := exec.Command("docker", "run", "-d", "-w", "/go/src/pcfdev", "-v", pwd+":/go/src/pcfdev", "pcfdev/provision", "sleep", "1000").Output()
	Expect(err).NotTo(HaveOccurred())
	dockerID = strings.TrimSpace(string(output))

	Expect(exec.Command("bash", "-c", "echo \"#!/bin/bash\necho 'Waiting for services to start...'\necho \\$@\" > "+pwd+"/provision-script").Run()).To(Succeed())
	Expect(exec.Command("docker", "exec", dockerID, "chmod", "+x", "/go/src/pcfdev/provision-script").Run()).To(Succeed())
	Expect(exec.Command("docker", "exec", dockerID, "go", "build", "-ldflags", "-X main.provisionScriptPath=/go/src/pcfdev/provision-script", "pcfdev").Run()).To(Succeed())
})

var _ = AfterSuite(func() {
	os.RemoveAll(filepath.Join(pwd, "pcfdev"))
	os.RemoveAll(filepath.Join(pwd, "provision-script"))
	Expect(exec.Command("docker", "rm", dockerID, "-f").Run()).To(Succeed())
})

var _ = Describe("PCF Dev provision", func() {
	It("should provision PCF Dev", func() {
		session, err := gexec.Start(exec.Command("docker", "exec", dockerID, "/go/src/pcfdev/pcfdev", "local.pcfdev.io"), GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))
		Expect(session).To(gbytes.Say("Waiting for services to start..."))
		Expect(session).To(gbytes.Say("To begin using PCF Dev, please run:"))
		Expect(session).To(gbytes.Say("cf login -a https://api.local.pcfdev.io --skip-ssl-validation"))
	})

	It("should create certificates", func() {
		session, err := gexec.Start(exec.Command("docker", "exec", dockerID, "/go/src/pcfdev/pcfdev", "local.pcfdev.io"), GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))

		session, err = gexec.Start(exec.Command("docker", "exec", dockerID, "bash", "-c", "echo 127.0.0.1 local.pcfdev.io >> /etc/hosts"), GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))

		session, err = gexec.Start(exec.Command("docker", "exec", dockerID, "service", "nginx", "start"), GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))

		session, err = gexec.Start(exec.Command("docker", "exec", dockerID, "curl", "--cacert", "/var/pcfdev/openssl/ca_cert.pem", "https://local.pcfdev.io:443"), GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))
	})

	Context("when provisioning fails", func() {
		BeforeEach(func() {
			Expect(exec.Command("bash", "-c", "echo \"#!/bin/bash\nexit 42\" > "+pwd+"/provision-script").Run()).To(Succeed())
		})

		It("should exit with the exit status of the provision script", func() {
			session, _ := gexec.Start(exec.Command("docker", "exec", dockerID, "/go/src/pcfdev/pcfdev", "local.pcfdev.io"), GinkgoWriter, GinkgoWriter)
			Eventually(session).Should(gexec.Exit(42))
		})
	})
})
