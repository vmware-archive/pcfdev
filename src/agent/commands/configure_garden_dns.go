package commands

import (
	"fmt"
	"agent/provisioner"
	"regexp"
	"strings"
	"agent/fs"
)

type ConfigureGardenDNS struct {
	FS        FS
	CmdRunner CmdRunner
}

func (c *ConfigureGardenDNS) Run() error {
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

	gardenBytes, err := c.FS.Read("/var/vcap/jobs/garden/bin/garden_ctl")
	if err != nil {
		return err
	}

	cleanedGardenCtl := []string{}

	for _, line := range strings.Split(string(gardenBytes), "\n") {
		if !strings.Contains(line, "-dnsServer=") {
			cleanedGardenCtl = append(cleanedGardenCtl, line)
		}
	}

	dnsInsertString := strings.Replace(strings.Join(cleanedGardenCtl, "\n"), `1>>$LOG_DIR/garden.stdout.log \`, fmt.Sprintf("-dnsServer=%s \\\n  1>>$LOG_DIR/garden.stdout.log \\", internalIP), fs.FileModeRootReadWrite)
	return c.FS.Write("/var/vcap/jobs/garden/bin/garden_ctl", strings.NewReader(dnsInsertString), fs.FileModeRootReadWrite)
}

func (*ConfigureGardenDNS) Distro() string {
	return provisioner.DistributionOSS
}
