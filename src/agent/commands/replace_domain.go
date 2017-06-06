package commands

import (
	"fmt"
	"agent/fs"
	"agent/provisioner"
	"strings"
)

type ReplaceDomain struct {
	CmdRunner provisioner.CmdRunner
	FS        provisioner.FS
	NewDomain string
}

func (r *ReplaceDomain) Run() error {
	files, err := r.CmdRunner.Output("bash", "-c", "find /var/vcap/jobs/*/ -type f")
	if err != nil {
		return err
	}

	oldDomain, err := r.FS.Read("/var/pcfdev/domain")
	if err != nil {
		return err
	}

	err = r.CmdRunner.Run("bash", "-c", fmt.Sprintf(`perl -p -i -e s/\\Q%s\\E/%s/g %s`, strings.TrimSpace(string(oldDomain)), r.NewDomain, strings.Replace(string(files), "\n", " ", -1)))
	if err != nil {
		return err
	}

	return r.FS.Write("/var/pcfdev/domain", strings.NewReader(r.NewDomain), fs.FileModeRootReadWrite)
}

func (r *ReplaceDomain) Distro() string {
	return provisioner.DistributionOSS
}
