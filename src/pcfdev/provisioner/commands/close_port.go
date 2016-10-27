package commands

import (
	"pcfdev/provisioner"
	"strings"
)

type ClosePort struct {
	CmdRunner provisioner.CmdRunner
	Port      string
}

func (c *ClosePort) Run() error {
	checkOutput, err := c.CmdRunner.Output("iptables", "-C", "INPUT", "-i", "eth1", "-p", "tcp", "--dport", c.Port, "-j", "REJECT")
	if err == nil {
		return nil
	}

	if strings.Contains(string(checkOutput), "No chain/target/match by that name") {
		return c.CmdRunner.Run("iptables", "-I", "INPUT", "-i", "eth1", "-p", "tcp", "--dport", c.Port, "-j", "REJECT")
	} else {
		return err
	}
}

func (*ClosePort) Distro() string {
	return provisioner.DistributionOSS
}