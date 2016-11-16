package commands_test

import (
	"errors"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"provisioner/provisioner"
	"provisioner/provisioner/commands"
	"provisioner/provisioner/mocks"
	"strings"
	"os"
	"provisioner/fs"
)

var _ = Describe("AddTrustedCerts", func() {
	var (
		mockCtrl      *gomock.Controller
		mockFS        = mocks.NewMockFS(mockCtrl)
		cmd           *commands.AddTrustedCerts
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockFS = mocks.NewMockCmdRunner(mockCtrl)
		cmd = &commands.AddTrustedCerts{
			FS: mockFS,
		}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("#Run", func() {
		It("should add an iptables rule to reject external access to the port", func() {
			gomock.InOrder(
				mockFS.EXPECT().Exists("/var/pcfdev/trusted_ca.crt").Return(true, nil),
				mockFS.EXPECT().Read("/var/pcfdev/trusted_ca.crt").Return([]byte("some-trusted-ca\n"), nil),
				mockFS.EXPECT().Exists("/var/vcap/jobs/gorouter/config/cert.pem").Return(true, nil),
				mockFS.EXPECT().Read("/var/vcap/jobs/gorouter/config/cert.pem").Return([]byte("some-gorouter-cert\n"), nil),
				mockFS.EXPECT().Write(
					"/var/vcap/jobs/cflinuxfs2-rootfs-setup/config/certs/trusted_ca.crt",
					strings.NewReader("some-trusted-ca\nsome-gorouter-cert\n"),
					os.FileMode(fs.FileModeRootReadWrite),
				).Return(nil),
			)

			Expect(cmd.Run()).To(Succeed())
		})

	})

	Describe("#Distro", func() {
		It("should return 'oss", func() {
			Expect(cmd.Distro()).To(Equal(provisioner.DistributionOSS))
		})
	})
})
