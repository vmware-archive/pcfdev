package provisioner

import (
	"bytes"
	"io"
)

//go:generate mockgen -package mocks -destination mocks/cert.go pcfdev/provisioner Cert
type Cert interface {
	GenerateCerts(domain string) (certificate []byte, privateKey []byte, caCertificate []byte, caPrivateKey []byte, err error)
}

//go:generate mockgen -package mocks -destination mocks/cmd_runner.go pcfdev/provisioner CmdRunner
type CmdRunner interface {
	Run(command string, args ...string) error
	Output(command string, args ...string) (output []byte, err error)
}

//go:generate mockgen -package mocks -destination mocks/fs.go pcfdev/provisioner FS
type FS interface {
	Mkdir(directory string) error
	Write(path string, contents io.Reader) error
	Read(path string) (contents []byte, err error)
	Exists(path string) (bool, error)
}

//go:generate mockgen -package mocks -destination mocks/ui.go pcfdev/provisioner UI
type UI interface {
	PrintHelpText(domain string) error
}

//go:generate mockgen -package mocks -destination mocks/command.go pcfdev/provisioner Command
type Command interface {
	Run() error
	Distro() string
}

type Provisioner struct {
	Cert             Cert
	CmdRunner        CmdRunner
	FS               FS
	UI               UI
	DisableUAAHSTS   Command
	ConfigureDnsmasq Command
	Commands         []Command

	Distro string
}

const (
	DistributionOSS = "oss"
	DistributionPCF = "pcf"
)

func (p *Provisioner) Provision(provisionScriptPath string, args ...string) error {
	domain := args[0]

	cert, key, caCert, _, err := p.Cert.GenerateCerts(domain)
	if err != nil {
		return err
	}

	if err := p.FS.Mkdir("/var/vcap/jobs/gorouter/config"); err != nil {
		return err
	}

	if err := p.FS.Write("/var/vcap/jobs/gorouter/config/cert.pem", bytes.NewReader(cert)); err != nil {
		return err
	}

	if err := p.FS.Write("/var/vcap/jobs/gorouter/config/key.pem", bytes.NewReader(key)); err != nil {
		return err
	}

	if err := p.FS.Mkdir("/var/pcfdev/openssl"); err != nil {
		return err
	}

	if err := p.FS.Write("/var/pcfdev/openssl/ca_cert.pem", bytes.NewReader(caCert)); err != nil {
		return err
	}

	for _, command := range p.Commands {
		if p.Distro == DistributionOSS && command.Distro() == DistributionPCF {
			continue
		}
		if err := command.Run(); err != nil {
			return err
		}
	}

	if err := p.CmdRunner.Run(provisionScriptPath, args...); err != nil {
		return err
	}

	return p.FS.Write("/run/pcfdev-healthcheck", bytes.NewReader([]byte("")))
}
