package commands_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/golang/mock/gomock"

	"agent/fs"
	"agent/provisioner"
	"agent/mocks"
	"agent/commands"
	"strings"
	"os"
)

var _ = Describe("SetupCFDot", func() {
	var (
		mockCtrl      *gomock.Controller
		mockFS        *mocks.MockFS
		mockCmdRunner *mocks.MockCmdRunner
		cmd           *commands.SetupCFDot
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockFS = mocks.NewMockFS(mockCtrl)
		mockCmdRunner = mocks.NewMockCmdRunner(mockCtrl)
		cmd = &commands.SetupCFDot{
			FS:        mockFS,
			CmdRunner: mockCmdRunner,
		}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("#Run", func() {
		It("should create a new file with content copied from /var/vcap/jobs/cfdot/bin/setup", func() {
			mockFS.EXPECT().Read("/var/vcap/jobs/cfdot/bin/setup").Return([]byte("cf-dot-setup-stuff-here\n"), nil)
			mockFS.EXPECT().Write("/etc/profile.d/cfdot.sh", strings.NewReader("cf-dot-setup-stuff-here\n"), os.FileMode(fs.FileModeRootReadWrite))

			Expect(cmd.Run()).To(Succeed())
		})
	})

	Describe("#Distro", func() {
		It("should return 'oss'", func() {
			Expect(cmd.Distro()).To(Equal(provisioner.DistributionOSS))
		})
	})
})
