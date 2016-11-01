package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"pcfdev/cert"
	"pcfdev/fs"
	"pcfdev/provisioner"
	"pcfdev/provisioner/commands"
	"strconv"
	"syscall"
	"time"
)

var (
	provisionScriptPath = "/var/pcfdev/run"
	timeoutInSeconds    = "3600"
	distro              = "pcf"
)
const (
	mysqlPort = "4568"
	rabbitBrokerPort = "4567"
	rabbitClusterDaemonPort = "25672"
)

func main() {
	provisionTimeout, err := strconv.Atoi(timeoutInSeconds)
	if err != nil {
		fmt.Printf("Error: %s.", err)
		os.Exit(1)
	}

	silentCommandRunner := &provisioner.ConcreteCmdRunner{
		Stdout:  ioutil.Discard,
		Stderr:  ioutil.Discard,
		Timeout: time.Duration(provisionTimeout) * time.Second,
	}
	p := &provisioner.Provisioner{
		Cert: &cert.Cert{},
		CmdRunner: &provisioner.ConcreteCmdRunner{
			Stdout:  os.Stdout,
			Stderr:  os.Stderr,
			Timeout: time.Duration(provisionTimeout) * time.Second,
		},
		FS: &fs.FS{},
		Commands: []provisioner.Command{
			&commands.DisableUAAHSTS{
				WebXMLPath: "/var/vcap/packages/uaa/tomcat/conf/web.xml",
			},
			&commands.ConfigureDnsmasq{
				Domain:     os.Args[1],
				ExternalIP: os.Args[2],
				FS:         &fs.FS{},
				CmdRunner:  silentCommandRunner,
			},
			&commands.ConfigureGardenDNS{
				FS:        &fs.FS{},
				CmdRunner: silentCommandRunner,
			},
			&commands.ClosePort{
				CmdRunner: silentCommandRunner,
				Port:      mysqlPort,
			},
			&commands.ClosePort{
				CmdRunner: silentCommandRunner,
				Port:      rabbitBrokerPort,
			},
			&commands.ClosePort{
				CmdRunner: silentCommandRunner,
				Port:      rabbitClusterDaemonPort,
			},
			&commands.SetupApi{
				CmdRunner: silentCommandRunner,
				FS:	&fs.FS{},
			},
		},

		Distro: distro,
	}

	if err := p.Provision(provisionScriptPath, os.Args[1:]...); err != nil {
		switch err.(type) {
		case *exec.ExitError:
			if exitErr, ok := err.(*exec.ExitError); ok {
				if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
					os.Exit(status.ExitStatus())
				} else {
					os.Exit(1)
				}
			}
		case *provisioner.TimeoutError:
			fmt.Printf("Timed out after %s seconds.\n", timeoutInSeconds)
			os.Exit(1)
		default:
			os.Exit(1)
		}
	}
}
