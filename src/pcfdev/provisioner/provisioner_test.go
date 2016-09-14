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
			mockUI        *mocks.MockUI
		)

		BeforeEach(func() {
			mockCtrl = gomock.NewController(GinkgoT())
			mockCmdRunner = mocks.NewMockCmdRunner(mockCtrl)
			mockUI = mocks.NewMockUI(mockCtrl)

			p = &provisioner.Provisioner{
				CmdRunner: mockCmdRunner,
				UI:        mockUI,
			}
		})

		AfterEach(func() {
			mockCtrl.Finish()
		})

		It("should provision a VM", func() {
			gomock.InOrder(
				mockCmdRunner.EXPECT().Run("some-provision-script-path", "some-domain"),
				mockUI.EXPECT().PrintHelpText("some-domain"),
			)

			Expect(p.Provision("some-provision-script-path", "some-domain")).To(Succeed())
		})

		Context("when there is an error running the provision script", func() {
			It("should return the error", func() {
				mockCmdRunner.EXPECT().Run("some-provision-script-path", "some-domain").Return(errors.New("some-error"))

				Expect(p.Provision("some-provision-script-path", "some-domain")).To(MatchError("some-error"))
			})
		})

		Context("when there is an error printing help text", func() {
			It("should return the error", func() {
				gomock.InOrder(
					mockCmdRunner.EXPECT().Run("some-provision-script-path", "some-domain"),
					mockUI.EXPECT().PrintHelpText("some-domain").Return(errors.New("some-error")),
				)

				Expect(p.Provision("some-provision-script-path", "some-domain")).To(MatchError("some-error"))
			})
		})
	})
})
