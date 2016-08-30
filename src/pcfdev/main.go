package main

import (
	"os"
	"os/exec"
	"syscall"
)

var provisionScriptPath = "/var/pcfdev/run"

func main() {
	cmd := exec.Command(provisionScriptPath, os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
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
