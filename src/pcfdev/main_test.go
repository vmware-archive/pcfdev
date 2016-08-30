package main_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var (
	binaryPath string
)

var _ = BeforeSuite(func() {
	tempDir, err := ioutil.TempDir("", "")
	Expect(err).NotTo(HaveOccurred())

	Expect(
		ioutil.WriteFile(
			filepath.Join(tempDir, "provision_script"),
			[]byte("#!/bin/bash\necho 'Waiting for services to start...'\necho $@"),
			0755),
	).To(Succeed())

	binaryPath, err = gexec.Build(
		"pcfdev",
		"-ldflags",
		"-X main.provisionScriptPath="+filepath.Join(tempDir, "provision_script"),
	)
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	os.RemoveAll(binaryPath)
})

var _ = Describe("PCF Dev provision", func() {
	It("should provision PCF Dev", func() {
		session, err := gexec.Start(exec.Command(binaryPath), GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))
		Expect(session).To(gbytes.Say("Waiting for services to start..."))
	})

	It("should pass arguments along", func() {
		session, err := gexec.Start(exec.Command(binaryPath, "local.pcfdev.io"), GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))
		Expect(session).To(gbytes.Say("Waiting for services to start..."))
		Expect(session).To(gbytes.Say("local.pcfdev.io"))
	})

	Context("when provisioning fails", func() {
		var failingBinaryPath string

		BeforeEach(func() {
			tempDir, err := ioutil.TempDir("", "")
			Expect(err).NotTo(HaveOccurred())

			Expect(
				ioutil.WriteFile(
					filepath.Join(tempDir, "provision_script"),
					[]byte("#!/bin/bash\nexit 42"),
					0755),
			).To(Succeed())

			failingBinaryPath, err = gexec.Build(
				"pcfdev",
				"-ldflags",
				"-X main.provisionScriptPath="+filepath.Join(tempDir, "provision_script"),
			)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should exit with the exit status of the provision script", func() {
			session, _ := gexec.Start(exec.Command(failingBinaryPath), GinkgoWriter, GinkgoWriter)
			Eventually(session).Should(gexec.Exit(42))
		})
	})
})
