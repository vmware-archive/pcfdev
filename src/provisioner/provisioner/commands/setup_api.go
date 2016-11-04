package commands

import (
	"provisioner/provisioner"
	"strings"
)

type SetupApi struct {
	CmdRunner provisioner.CmdRunner
	FS	provisioner.FS
}

func (s *SetupApi) Run() error {
	monitrcContents := `check process pcfdev-api
  with pidfile /var/pcfdev/api/api.pid
  start program "/var/pcfdev/api/api_ctl start"
  stop program "/var/pcfdev/api/api_ctl stop"
  group vcap
  mode manual`

	apiCtlContents := `#!/bin/bash
set -ex

case $1 in

  PIDFILE=/var/vcap/sys/run/pcfdev-api/api.pid

  start)
    mkdir -p /var/vcap/sys/run/pcfdev-api
    sleep 999999999 &
    echo $! > ${PIDFILE}

    ;;

  stop)
    kill $(cat $PIDFILE)

    ;;

  *)
    echo "Usage: pcfdev_api_ctl {start|stop}"
    ;;

esac`

	err := s.FS.Write("/var/vcap/monit/job/1001_pcfdev_api.monitrc", strings.NewReader(monitrcContents))
	if err != nil {
		return err
	}

	if err := s.FS.Mkdir("/var/pcfdev/api"); err !=nil {
		return err
	}

	return s.FS.Write("/var/pcfdev/api/api_ctl", strings.NewReader(apiCtlContents))
}

func(s *SetupApi) Distro() string {
	return provisioner.DistributionOSS
}