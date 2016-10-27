package commands_test

import (
	"errors"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"pcfdev/provisioner"
	"pcfdev/provisioner/commands"
	"pcfdev/provisioner/mocks"
)

var _ = Describe("ClosePort", func() {
	var (
		mockCtrl      *gomock.Controller
		mockCmdRunner *mocks.MockCmdRunner
		cmd           *commands.ClosePort
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockCmdRunner = mocks.NewMockCmdRunner(mockCtrl)
		cmd = &commands.ClosePort{
			CmdRunner: mockCmdRunner,
			Port:      "some-port",
		}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("#Run", func() {
		It("should add an iptables rule to reject external access to the port", func() {
			gomock.InOrder(
				mockCmdRunner.EXPECT().
					Output("iptables", "-C", "INPUT", "-i", "eth1", "-p", "tcp", "--dport", "some-port", "-j", "REJECT").
					Return([]byte("iptables: No chain/target/match by that name."), errors.New("error-because-rule-does-not-exist")),
				mockCmdRunner.EXPECT().Run("iptables", "-I", "INPUT", "-i", "eth1", "-p", "tcp", "--dport", "some-port", "-j", "REJECT"),
			)

			Expect(cmd.Run()).To(Succeed())
		})

		Context("when the iptables rule already exists", func() {
			It("should not add a new rule", func() {
				mockCmdRunner.EXPECT().
					Output("iptables", "-C", "INPUT", "-i", "eth1", "-p", "tcp", "--dport", "some-port", "-j", "REJECT").
					Return(nil, nil)

				Expect(cmd.Run()).To(Succeed())
			})
		})

		Context("when there is an error checking the iptables rule", func() {
			It("should return the error", func() {
				mockCmdRunner.EXPECT().
					Output("iptables", "-C", "INPUT", "-i", "eth1", "-p", "tcp", "--dport", "some-port", "-j", "REJECT").
					Return(nil, errors.New("some-error"))

				Expect(cmd.Run()).To(MatchError("some-error"))
			})
		})

		Context("when there is an error inserting the iptables rule", func() {
			It("should return the error", func() {
				gomock.InOrder(
					mockCmdRunner.EXPECT().
						Output("iptables", "-C", "INPUT", "-i", "eth1", "-p", "tcp", "--dport", "some-port", "-j", "REJECT").
						Return([]byte("iptables: No chain/target/match by that name."), errors.New("error-because-rule-does-not-exist")),
					mockCmdRunner.EXPECT().
						Run("iptables", "-I", "INPUT", "-i", "eth1", "-p", "tcp", "--dport", "some-port", "-j", "REJECT").
						Return(errors.New("some-error")),
				)

				Expect(cmd.Run()).To(MatchError("some-error"))
			})
		})
	})

	Describe("#Distro", func() {
		It("should return 'oss", func() {
			Expect(cmd.Distro()).To(Equal(provisioner.DistributionOSS))
		})
	})
})
