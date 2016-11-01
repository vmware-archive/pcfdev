package commands_test

import (
	"github.com/golang/mock/gomock"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"pcfdev/provisioner/mocks"
	"pcfdev/provisioner/commands"
	"pcfdev/provisioner"
	"errors"
)

var _ = Describe("SetupApi", func() {
	var (
		mockCtrl      *gomock.Controller
		mockFS        *mocks.MockFS
		mockCmdRunner *mocks.MockCmdRunner
		cmd           *commands.SetupApi
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockFS = mocks.NewMockFS(mockCtrl)
		mockCmdRunner = mocks.NewMockCmdRunner(mockCtrl)
		cmd = &commands.SetupApi{
			FS:        mockFS,
			CmdRunner: mockCmdRunner,
		}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("#Run", func() {
		Context("When the file system is in a bad state", func() {
			It("returns the error from failing to write the monitrc", func() {
				mockFS.EXPECT().Write("/var/pcfdev/api/api_ctl", gomock.Any()).AnyTimes()
				mockFS.EXPECT().Write("/var/vcap/monit/job/1001_pcfdev_api.monitrc", gomock.Any()).Return(errors.New("some-error"))

				Expect(cmd.Run()).To(MatchError("some-error"))
			})

			It("returns the error from failing to write the api_ctl", func() {
				mockFS.EXPECT().Write("/var/vcap/monit/job/1001_pcfdev_api.monitrc", gomock.Any()).AnyTimes()
				mockFS.EXPECT().Write("/var/pcfdev/api/api_ctl", gomock.Any()).Return(errors.New("some-error"))

				Expect(cmd.Run()).To(MatchError("some-error"))
			})
		})

		It("write a monit file to the /var/vcap/monit/job", func() {
			monitrc := `check process pcfdev-api
  with pidfile /var/pcfdev/api/api.pid
  start program "/var/pcfdev/api/api_ctl start"
  stop program "/var/pcfdev/api/api_ctl stop"
  group vcap
  mode manual`

			monit_ctl := `#!/bin/bash
set -ex

case $1 in

  PIDFILE=/var/pcfdev/api/api.pid

  start)
    /var/pcfdev/api/api &
    echo $! > ${PIDFILE}

    ;;

  stop)
    kill $(cat $PIDFILE)

    ;;

  *)
    echo "Usage: pcfdev_api_ctl {start|stop}"
    ;;

esac`

			gomock.InOrder(
				mockFS.EXPECT().Write("/var/vcap/monit/job/1001_pcfdev_api.monitrc", strings.NewReader(monitrc)),
				mockFS.EXPECT().Write("/var/pcfdev/api/api_ctl", strings.NewReader(monit_ctl)),
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