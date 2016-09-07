package provisioner

import (
	"io"
	"os/exec"
)

type ConcreteCmdRunner struct {
	Stdout io.Writer
	Stderr io.Writer
}

func (r *ConcreteCmdRunner) Run(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = r.Stdout
	cmd.Stderr = r.Stderr
	return cmd.Run()
}
