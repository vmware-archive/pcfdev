package provisioner

import (
	"fmt"
	"io"
	"os/exec"
	"syscall"
	"time"
)

type ConcreteCmdRunner struct {
	Stdout  io.Writer
	Stderr  io.Writer
	Timeout time.Duration
}

func (r *ConcreteCmdRunner) Run(command string, args ...string) error {
	fmt.Fprintln(r.Stdout, fmt.Sprintf("timeout duration is %+v\n", r.Timeout))

	cmd := exec.Command(command, args...)
	cmd.Stdout = r.Stdout
	cmd.Stderr = r.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		return err
	}

	timer := time.AfterFunc(r.Timeout, func() {
		pgid, err := syscall.Getpgid(cmd.Process.Pid)
		if err == nil {
			syscall.Kill(-pgid, 15)
		}
	})

	err := cmd.Wait()

	if !timer.Stop() {
		return &timeoutError{}
	}

	return err
}
