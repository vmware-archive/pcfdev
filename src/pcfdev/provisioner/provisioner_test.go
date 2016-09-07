package provisioner_test

import (
	"errors"
	"pcfdev/provisioner"
	"pcfdev/provisioner/mocks"

	"github.com/golang/mock/gomock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Provisioner", func() {
	Describe("#Provision", func() {
		var (
			p             *provisioner.Provisioner
			mockCtrl      *gomock.Controller
			mockCmdRunner *mocks.MockCmdRunner
		)

		BeforeEach(func() {
			mockCtrl = gomock.NewController(GinkgoT())
			mockCmdRunner = mocks.NewMockCmdRunner(mockCtrl)

			p = &provisioner.Provisioner{
				CmdRunner: mockCmdRunner,
			}
		})

		AfterEach(func() {
			mockCtrl.Finish()
		})

		It("should provision a VM", func() {
			mockCmdRunner.EXPECT().Run("some-path", "some-arg")

			Expect(p.Provision("some-path", "some-arg")).To(Succeed())
		})

		Context("when there is an error running the provision script", func() {
			It("should return the error", func() {
				mockCmdRunner.EXPECT().Run("some-path", "some-arg").Return(errors.New("some-error"))

				Expect(p.Provision("some-path", "some-arg")).To(MatchError("some-error"))
			})
		})
	})
})
