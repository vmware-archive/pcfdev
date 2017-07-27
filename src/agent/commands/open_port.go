package commands

import (
	"agent/provisioner"
)

type OpenPort struct {
	CmdRunner CmdRunner
	Port string
}

func (o *OpenPort)  Run() error {
	return o.CmdRunner.Run("iptables", "-I", "INPUT", "-p", "tcp", "--dport", o.Port, "-j", "ACCEPT")
}

func (*OpenPort) Distro() string {
	return provisioner.DistributionOSS
}