package commands_test

import (
	"provisioner/provisioner/commands"
	"provisioner/provisioner/mocks"
	"strings"

	"github.com/cppforlife/packer-bosh/bosh-provisioner/Godeps/_workspace/src/github.com/cloudfoundry/bosh-agent/errors"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"provisioner/provisioner"
	"os"
	"provisioner/fs"
)

var _ = Describe("ConfigureGardenDNS", func() {
	var (
		mockCtrl      *gomock.Controller
		mockFS        *mocks.MockFS
		mockCmdRunner *mocks.MockCmdRunner
		cmd           *commands.ConfigureGardenDNS
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockFS = mocks.NewMockFS(mockCtrl)
		mockCmdRunner = mocks.NewMockCmdRunner(mockCtrl)
		cmd = &commands.ConfigureGardenDNS{
			FS:        mockFS,
			CmdRunner: mockCmdRunner,
		}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("#Run", func() {
		Context("when there are no dnsServers", func() {
			It("should use the internal IP as a DNS server in garden", func() {
				gomock.InOrder(
					mockCmdRunner.EXPECT().Output("ip", "route", "get", "1").Return([]byte("some-ip via some-other-ip dev eth0  src some-internal-ip\n cache"), nil),
					mockFS.EXPECT().Read("/var/vcap/jobs/garden/bin/garden_ctl").Return([]byte("some-executable \\\n  -someProperty=some-value \\\n  1>>$LOG_DIR/garden.stdout.log \\"), nil),
					mockFS.EXPECT().Write("/var/vcap/jobs/garden/bin/garden_ctl", strings.NewReader("some-executable \\\n  -someProperty=some-value \\\n  -dnsServer=some-internal-ip \\\n  1>>$LOG_DIR/garden.stdout.log \\"), os.FileMode(fs.FileModeRootReadWrite)),
				)

				Expect(cmd.Run()).To(Succeed())
			})
		})

		Context("when there are dnsServers", func() {
			It("should use the internal IP as a DNS server in garden", func() {
				gomock.InOrder(
					mockCmdRunner.EXPECT().Output("ip", "route", "get", "1").Return([]byte("some-ip via some-other-ip dev eth0  src some-internal-ip\n cache"), nil),
					mockFS.EXPECT().Read("/var/vcap/jobs/garden/bin/garden_ctl").Return([]byte("some-executable \\\n  -someProperty=some-value \\\n  -dnsServer=some-ip \\\n  -dnsServer=some-other-ip \\\n  1>>$LOG_DIR/garden.stdout.log \\"), nil),
					mockFS.EXPECT().Write("/var/vcap/jobs/garden/bin/garden_ctl", strings.NewReader("some-executable \\\n  -someProperty=some-value \\\n  -dnsServer=some-internal-ip \\\n  1>>$LOG_DIR/garden.stdout.log \\"), os.FileMode(fs.FileModeRootReadWrite)),
				)

				Expect(cmd.Run()).To(Succeed())
			})
		})

		Context("when there is an error gettting the internal ip", func() {
			It("return the error", func() {
				mockCmdRunner.EXPECT().Output("ip", "route", "get", "1").Return(nil, errors.New("some-error"))

				Expect(cmd.Run()).To(MatchError("some-error"))
			})
		})

		Context("when the internal ip cannot be parsed", func() {
			It("return an error", func() {
				mockCmdRunner.EXPECT().Output("ip", "route", "get", "1").Return([]byte("some-bad-output"), nil)

				Expect(cmd.Run()).To(MatchError("internal ip could not be parsed from output: some-bad-output"))
			})
		})

		Context("when there is an error reading the garden ctl", func() {
			It("return the error", func() {
				gomock.InOrder(
					mockCmdRunner.EXPECT().Output("ip", "route", "get", "1").Return([]byte("some-ip via some-other-ip dev eth0  src some-internal-ip\n cache"), nil),
					mockFS.EXPECT().Read("/var/vcap/jobs/garden/bin/garden_ctl").Return(nil, errors.New("some-error")),
				)

				Expect(cmd.Run()).To(MatchError("some-error"))
			})
		})

		Context("when there is an error rewriting the garden ctl", func() {
			It("return the error", func() {
				gomock.InOrder(
					mockCmdRunner.EXPECT().Output("ip", "route", "get", "1").Return([]byte("some-ip via some-other-ip dev eth0  src some-internal-ip\n cache"), nil),
					mockFS.EXPECT().Read("/var/vcap/jobs/garden/bin/garden_ctl").Return([]byte("some-executable \\\n  -someProperty=some-value \\\n  -dnsServer=some-ip \\\n  -dnsServer=some-other-ip \\\n  1>>$LOG_DIR/garden.stdout.log \\"), nil),
					mockFS.EXPECT().Write("/var/vcap/jobs/garden/bin/garden_ctl", strings.NewReader("some-executable \\\n  -someProperty=some-value \\\n  -dnsServer=some-internal-ip \\\n  1>>$LOG_DIR/garden.stdout.log \\"), os.FileMode(fs.FileModeRootReadWrite)).Return(errors.New("some-error")),
				)

				Expect(cmd.Run()).To(MatchError("some-error"))
			})
		})
	})

	Describe("#Distro", func() {
		It("should return 'oss'", func() {
			Expect(cmd.Distro()).To(Equal(provisioner.DistributionOSS))
		})
	})
})
