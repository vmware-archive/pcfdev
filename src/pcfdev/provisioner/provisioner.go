package provisioner

//go:generate mockgen -package mocks -destination mocks/cmd_runner.go pcfdev/provisioner CmdRunner
type CmdRunner interface {
	Run(command string, args ...string) error
}

//go:generate mockgen -package mocks -destination mocks/ui.go pcfdev/provisioner UI
type UI interface {
	PrintHelpText(domain string) error
}

type Provisioner struct {
	CmdRunner CmdRunner
	UI        UI
}

func (p *Provisioner) Provision(provisionScriptPath string, args ...string) error {
	if err := p.CmdRunner.Run(provisionScriptPath, args...); err != nil {
		return err
	}
	return p.UI.PrintHelpText(args[0])
}
