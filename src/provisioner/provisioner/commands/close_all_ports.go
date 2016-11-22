package commands

import (
	"provisioner/provisioner"
)

type CloseAllPorts struct {
	CmdRunner provisioner.CmdRunner
}

func (c *CloseAllPorts) Run() error {
	if err := c.dropNewConnections("eth0"); err != nil {
		return err
	}

	if err := c.dropNewConnections("eth1"); err != nil {
		return err
	}

	return c.CmdRunner.Run("iptables", "-I", "INPUT", "-i", "lo", "-j", "ACCEPT")
}

func (c *CloseAllPorts) dropNewConnections(interfaceName string) error {
	if err := c.CmdRunner.Run("iptables", "-I", "INPUT", "-i", interfaceName, "-p", "tcp", "-j", "DROP"); err != nil {
		return err
	}

	return c.CmdRunner.Run("iptables", "-I", "INPUT", "-i", interfaceName, "-m", "conntrack", "--ctstate", "ESTABLISHED,RELATED", "-j", "ACCEPT")
}

func (*CloseAllPorts) Distro() string {
	return provisioner.DistributionOSS
}