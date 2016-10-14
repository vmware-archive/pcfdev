package commands

import (
	"fmt"
	"pcfdev/provisioner"
	"regexp"
	"strings"
)

type ConfigureDnsmasq struct {
	FS         provisioner.FS
	CmdRunner  provisioner.CmdRunner
	Domain     string
	ExternalIP string
}

func (c *ConfigureDnsmasq) Run() error {
	if err := c.CmdRunner.Run("resolvconf", "--disable-updates"); err != nil {
		return err
	}

	if err := c.CmdRunner.Run("service", "dnsmasq", "stop"); err != nil {
		return err
	}

	output, err := c.CmdRunner.Output("ip", "route", "get", "1")
	if err != nil {
		return err
	}

	var internalIP string
	regex := regexp.MustCompile(`\s{2}src\s(.*)`)
	if matches := regex.FindStringSubmatch(string(output)); len(matches) > 1 {
		internalIP = matches[1]
	} else {
		return fmt.Errorf("internal ip could not be parsed from output: %s", string(output))
	}

	if err := c.FS.Write("/etc/dnsmasq.d/domain", strings.NewReader(fmt.Sprintf("address=/.%s/%s\naddress=/.cf.internal/%s", c.Domain, c.ExternalIP, internalIP))); err != nil {
		return err
	}

	if err := c.FS.Write("/etc/dnsmasq.conf", strings.NewReader("resolv-file=/var/pcfdev/external-resolv.conf")); err != nil {
		return err
	}

	exists, err := c.FS.Exists("/var/pcfdev/external-resolv.conf")
	if err != nil {
		return err
	}

	if !exists {
		data, err := c.FS.Read("/etc/resolv.conf")
		if err != nil {
			return err
		}

		nameservers := []string{}
		for _, nameserver := range strings.Split(string(data), "\n") {
			if strings.HasPrefix(nameserver, "nameserver") && !strings.Contains(nameserver, "127.0.0.1") {
				nameservers = append(nameservers, nameserver)
			}
		}

		if err := c.FS.Write("/var/pcfdev/external-resolv.conf", strings.NewReader(strings.Join(nameservers, "\n"))); err != nil {
			return err
		}

	}

	if err := c.CmdRunner.Run("service", "dnsmasq", "start"); err != nil {
		return err
	}

	return c.FS.Write("/etc/resolv.conf", strings.NewReader(fmt.Sprintf("nameserver %s", internalIP)))
}
