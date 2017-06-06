package commands_test

import (
	"github.com/golang/mock/gomock"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"errors"
	"os"
	"agent/fs"
	"agent/provisioner"
	"agent/commands"
	"agent/mocks"
)

var _ = Describe("ReplaceDomain", func() {
	var (
		mockCtrl      *gomock.Controller
		mockFS        *mocks.MockFS
		mockCmdRunner *mocks.MockCmdRunner
		cmd           *commands.ReplaceDomain
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockFS = mocks.NewMockFS(mockCtrl)
		mockCmdRunner = mocks.NewMockCmdRunner(mockCtrl)
		cmd = &commands.ReplaceDomain{
			FS:        mockFS,
			CmdRunner: mockCmdRunner,
			NewDomain: "some-new-domain",
		}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("#Run", func() {
		It("should replace the old domain with the new domain in all files (without following symlinks) in /var/vcap/jobs", func() {
			gomock.InOrder(
				mockCmdRunner.EXPECT().Output("bash", "-c", "find /var/vcap/jobs/*/ -type f").Return([]byte("/var/vcap/jobs/some-job/some-file\n/var/vcap/jobs/some-job/some-other-file"), nil),
				mockFS.EXPECT().Read("/var/pcfdev/domain").Return([]byte("some-old-domain\n"), nil),
				mockCmdRunner.EXPECT().Run("bash", "-c", `perl -p -i -e s/\\Qsome-old-domain\\E/some-new-domain/g /var/vcap/jobs/some-job/some-file /var/vcap/jobs/some-job/some-other-file`),
				mockFS.EXPECT().Write("/var/pcfdev/domain", strings.NewReader("some-new-domain"), os.FileMode(fs.FileModeRootReadWrite)),
			)

			Expect(cmd.Run()).To(Succeed())
		})

		Context("when finding job files fails", func() {
			It("should return an error", func() {
				mockCmdRunner.EXPECT().Output("bash", "-c", "find /var/vcap/jobs/*/ -type f").Return(nil, errors.New("some-error"))

				Expect(cmd.Run()).To(MatchError("some-error"))
			})
		})

		Context("when reading domain file fails", func() {
			It("should return an error", func() {
				gomock.InOrder(
					mockCmdRunner.EXPECT().Output("bash", "-c", "find /var/vcap/jobs/*/ -type f").Return([]byte("/var/vcap/jobs/some-job/some-file\n/var/vcap/jobs/some-job/some-other-file"), nil),
					mockFS.EXPECT().Read("/var/pcfdev/domain").Return(nil, errors.New("some-error")),
				)
				Expect(cmd.Run()).To(MatchError("some-error"))
			})
		})

		Context("when replacing the domain in the job files fails", func() {
			It("should return an error", func() {
				gomock.InOrder(
					mockCmdRunner.EXPECT().Output("bash", "-c", "find /var/vcap/jobs/*/ -type f").Return([]byte("/var/vcap/jobs/some-job/some-file\n/var/vcap/jobs/some-job/some-other-file"), nil),
					mockFS.EXPECT().Read("/var/pcfdev/domain").Return([]byte("some-old-domain\n"), nil),
					mockCmdRunner.EXPECT().Run("bash", "-c", `perl -p -i -e s/\\Qsome-old-domain\\E/some-new-domain/g /var/vcap/jobs/some-job/some-file /var/vcap/jobs/some-job/some-other-file`).Return(errors.New("some-error")),
				)

				Expect(cmd.Run()).To(MatchError("some-error"))
			})
		})

		Context("when writing the new domain fails", func() {
			It("should return an error", func() {
				gomock.InOrder(
					mockCmdRunner.EXPECT().Output("bash", "-c", "find /var/vcap/jobs/*/ -type f").Return([]byte("/var/vcap/jobs/some-job/some-file\n/var/vcap/jobs/some-job/some-other-file"), nil),
					mockFS.EXPECT().Read("/var/pcfdev/domain").Return([]byte("some-old-domain\n"), nil),
					mockCmdRunner.EXPECT().Run("bash", "-c", `perl -p -i -e s/\\Qsome-old-domain\\E/some-new-domain/g /var/vcap/jobs/some-job/some-file /var/vcap/jobs/some-job/some-other-file`),
					mockFS.EXPECT().Write("/var/pcfdev/domain", strings.NewReader("some-new-domain"), os.FileMode(fs.FileModeRootReadWrite)).Return(errors.New("some-error")),
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
