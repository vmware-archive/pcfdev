package commands

import (
	"agent/provisioner"
	"strings"
	"agent/fs"
)

type SetupApi struct {
	CmdRunner CmdRunner
	FS	FS
}

func (s *SetupApi) Run() error {
	monitrcContents := `check process pcfdev-api
  with pidfile /var/vcap/sys/run/pcfdev-api/api.pid
  start program "/var/pcfdev/api/api_ctl start"
  stop program "/var/pcfdev/api/api_ctl stop"
  group vcap
  mode manual`

	apiCtlContents := `#!/bin/bash
set -ex

PIDFILE=/var/vcap/sys/run/pcfdev-api/api.pid

case $1 in

  start)
    mkdir -p /var/vcap/sys/run/pcfdev-api
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

	err := s.FS.Write("/var/vcap/monit/job/1001_pcfdev_api.monitrc", strings.NewReader(monitrcContents), fs.FileModeRootReadWrite)
	if err != nil {
		return err
	}

	return s.FS.Write("/var/pcfdev/api/api_ctl", strings.NewReader(apiCtlContents), fs.FileModeRootReadWriteExecutable)
}

func(s *SetupApi) Distro() string {
	return provisioner.DistributionOSS
}