package commands

import (
	"provisioner/provisioner"
	"strings"
)

type AddTrustedCerts struct {
	FS provisioner.FS
}

func (c *AddTrustedCerts) Run() error {
}

func (*AddTrustedCerts) Distro() string {
	return provisioner.DistributionOSS
}