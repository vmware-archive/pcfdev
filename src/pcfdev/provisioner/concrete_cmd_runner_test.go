package provisioner_test

import (
	"pcfdev/provisioner"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("ConcreteCmdRunner", func() {
	Describe("#Run", func() {
		var r *provisioner.ConcreteCmdRunner
		var stdout *gbytes.Buffer
		var stderr *gbytes.Buffer

		BeforeEach(func() {
			stdout = gbytes.NewBuffer()
			stderr = gbytes.NewBuffer()

			r = &provisioner.ConcreteCmdRunner{
				Stdout: stdout,
				Stderr: stderr,
			}
		})

		It("should run commands", func() {
			err := r.Run("echo", "-n", "some output")
			Expect(err).NotTo(HaveOccurred())
			Eventually(stdout).Should(gbytes.Say("some output"))

			err = r.Run("bash", "-c", ">&2 echo -n some output")
			Expect(err).NotTo(HaveOccurred())
			Eventually(stderr).Should(gbytes.Say("some output"))
		})

		Context("when there is an error", func() {
			It("should return the error and the output", func() {
				err := r.Run("/some/bad/binary")
				Expect(err).To(MatchError(ContainSubstring("no such file or directory")))
			})
		})
	})
})
