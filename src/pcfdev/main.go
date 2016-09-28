package main

import (
	"fmt"
	"os"
	"os/exec"
	"pcfdev/cert"
	"pcfdev/fs"
	"pcfdev/provisioner"
	"strconv"
	"syscall"
	"time"
)

var (
	provisionScriptPath = "/var/pcfdev/run"
	timeoutInSeconds    = "1800"
)

func main() {
	provisionTimeout, err := strconv.Atoi(timeoutInSeconds)
	if err != nil {
		fmt.Printf("Error: %s.", err)
		os.Exit(1)
	}

	p := &provisioner.Provisioner{
		Cert: &cert.Cert{},
		CmdRunner: &provisioner.ConcreteCmdRunner{
			Stdout:  os.Stdout,
			Stderr:  os.Stderr,
			Timeout: time.Duration(provisionTimeout) * time.Second,
		},
		FS: &fs.FS{},
	}

	if err := p.Provision(provisionScriptPath, os.Args[1:]...); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				os.Exit(status.ExitStatus())
			} else {
				os.Exit(1)
			}
		} else {
			fmt.Printf("Error: %s.", err)
			os.Exit(1)
		}
	}
}
