package main

import (
	"os"
	"os/exec"
	"pcfdev/provisioner"
	"syscall"
)

var provisionScriptPath = "/var/pcfdev/run"

func main() {
	p := &provisioner.Provisioner{
		CmdRunner: &provisioner.ConcreteCmdRunner{
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		},
	}

	if err := p.Provision(provisionScriptPath, os.Args[1:]...); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				os.Exit(status.ExitStatus())
			} else {
				os.Exit(1)
			}
		} else {
			os.Exit(1)
		}
	}
}
