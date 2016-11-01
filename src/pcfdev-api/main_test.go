package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"os"
	"os/exec"
	"strings"
	"path/filepath"
)

var _ = Describe("PCF Dev Api provision", func() {
	var (
		dockerID string
		pwd      string
	)

	BeforeEach(func() {
		var err error
		pwd, err = os.Getwd()
		Expect(err).NotTo(HaveOccurred())

		output, err := exec.Command("docker", "run", "--privileged", "-d", "-w", "/go/src/pcfdev-api", "-v", pwd+":/go/src/pcfdev-api", "pcfdev/provision", "bash", "-c", "mount -o size=40M -t tmpfs tmpfs /run && sleep 1000").Output()
		Expect(err).NotTo(HaveOccurred())
		dockerID = strings.TrimSpace(string(output))
		Expect(exec.Command("docker", "exec", dockerID, "go", "build", "pcfdev-api").Run()).To(Succeed())
		Expect(exec.Command("docker", "exec", dockerID, "bash", "-c", `/go/src/pcfdev-api/pcfdev-api`).Start())
		Eventually(func() error {
			return exec.Command("docker", "exec", dockerID, "lsof", "-i", ":8090").Run()
		}).Should(Succeed())
	})

	AfterEach(func() {
		os.RemoveAll(filepath.Join(pwd, "pcfdev"))
		os.RemoveAll(filepath.Join(pwd, "provisdion-script"))
		Expect(exec.Command("docker", "rm", dockerID, "-f").Run()).To(Succeed())
	})

	Describe("/replace-secrets", func() {
		It("should replace secrets on the VM", func() {
			session, err := gexec.Start(exec.Command("docker", "exec", dockerID, "curl", "-v", "-X", "PUT", "-d", `{"password":"some-master-password"}`, "localhost:8090/replace-secrets"), GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session, "10s").Should(gexec.Exit())

			session, err = gexec.Start(exec.Command("docker", "exec", dockerID, "cat", "/var/vcap/jobs/uaa/config/uaa.yml"), GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session, "10s").Should(gexec.Exit(0))
			Expect(session).To(gbytes.Say(`- admin\|some-master-password\|scim.write`))
		})
	})

	Describe("/status", func() {
		Context("when the health-check is not written", func() {
			It("should reply 'Unprovisioned' in the /status endpoint", func() {
				session, err := gexec.Start(exec.Command("docker", "exec", dockerID, "curl", "localhost:8090/status"), GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())
				Eventually(session, "10s").Should(gexec.Exit(0))
				Expect(session).To(gbytes.Say(`{"status":"Unprovisioned"}`))
			})
		})

		Context("when the health-check is written", func() {
			It("should reply `Running` in the /status endpoint", func() {
				session, err := gexec.Start(exec.Command("docker", "exec", dockerID, "touch", "/run/pcfdev-healthcheck"), GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())
				Eventually(session, "10s").Should(gexec.Exit(0))

				session, err = gexec.Start(exec.Command("docker", "exec", dockerID, "curl", "localhost:8090/status"), GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())
				Eventually(session, "10s").Should(gexec.Exit(0))
				Expect(session).To(gbytes.Say(`{"status":"Running"}`))
			})
		})
	})

})
