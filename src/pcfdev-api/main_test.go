package main_test

import (
	"os"
	"os/exec"
	"strings"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/gbytes"
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
	})

	AfterEach(func() {
		os.RemoveAll(filepath.Join(pwd, "pcfdev"))
		os.RemoveAll(filepath.Join(pwd, "provision-script"))
		Expect(exec.Command("docker", "rm", dockerID, "-f").Run()).To(Succeed())
	})


	It("should run PCF Dev Api sucessfully", func() {
		session, err := gexec.Start(exec.Command("docker", "exec", dockerID, "curl", "localhost:8090"), GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session, "10s").Should(gexec.Exit(0))
		Expect(session).To(gbytes.Say(`{"error":{"message":"ROUTE_NOT_FOUND"}}`))
	})

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
