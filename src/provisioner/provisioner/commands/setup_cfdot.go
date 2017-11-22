package commands

import (
	"provisioner/provisioner"
	"strings"
	"provisioner/fs"
)

type SetupCFDot struct {
	CmdRunner provisioner.CmdRunner
	FS        provisioner.FS
}

func (s *SetupCFDot) Run() error {
	setupFileContentsBytes, err := s.FS.Read("/var/vcap/jobs/cfdot/bin/setup")
	if err != nil {
		return err
	}

	return s.FS.Write("/etc/profile.d/cfdot.sh", strings.NewReader(string(setupFileContentsBytes)), fs.FileModeRootReadWrite)
}

func (s *SetupCFDot) Distro() string {
	return provisioner.DistributionOSS
}

