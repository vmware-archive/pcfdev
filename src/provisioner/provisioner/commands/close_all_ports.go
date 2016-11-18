package commands

import (
	"provisioner/provisioner"
)

type CloseAllPorts struct {
	CmdRunner provisioner.CmdRunner
}

func (c *CloseAllPorts) Run() error {
	err := c.CmdRunner.Run("iptables", "-I", "INPUT", "-p", "tcp", "-j", "DROP")
	if err != nil {
		return err
	}
	return c.CmdRunner.Run("iptables", "-I", "INPUT", "-i", "lo", "-j", "ACCEPT")
}

func (*CloseAllPorts) Distro() string {
	return provisioner.DistributionOSS
}