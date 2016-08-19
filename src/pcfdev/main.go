package main

import (
	"os"
	"os/exec"
)

var provisionScriptPath = "/var/pcfdev/run"

func main() {
	cmd := exec.Command(provisionScriptPath, os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
