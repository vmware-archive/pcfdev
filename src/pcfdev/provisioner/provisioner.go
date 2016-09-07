package provisioner

//go:generate mockgen -package mocks -destination mocks/cmd_runner.go pcfdev/provisioner CmdRunner
type CmdRunner interface {
	Run(command string, args ...string) error
}

type Provisioner struct {
	CmdRunner CmdRunner
}

func (p *Provisioner) Provision(provisionScriptPath string, args ...string) error {
	return p.CmdRunner.Run(provisionScriptPath, args...)
}
