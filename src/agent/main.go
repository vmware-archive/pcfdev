package main

import (
	"agent/cert"
	"agent/commands"
	"agent/fs"
	"agent/provisioner"
	"agent/runner"
	"agent/usecases"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"
	"errors"
)

var (
	provisionScriptPath = "/var/pcfdev/run"
	timeoutInSeconds    = "3600"
	distro              = "pcf"
)

type provisionForm struct {
	Domain           string
	IP               string
	Services         string
	DockerRegistries string
	Provider         string
}

func fileExists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func serverError(w http.ResponseWriter) {
	errorHandler(w, "Failed to replace UAA Config Credentials", http.StatusInternalServerError)
}

func errorHandler(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, fmt.Sprintf(`{"error":{"message":"%s"}}`, message))
}

func checkArgCount(form provisionForm) error {
	error := errors.New("Need 5 arguments, Usage: ./provision <domain> <ip> <services> <docker_registries> <provider>")

	if form.Domain == "" {
		return error
	}

	if form.IP == "" {
		return error
	}

	if form.Services == "" {
		return error
	}

	if form.DockerRegistries == "" {
		return error
	}
	
	if form.Provider == "" {
		return error
	}

	return nil
}

func provision(w http.ResponseWriter, r *http.Request) {
	var form provisionForm
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		http.Error(w, "Need 5 arguments, Usage: ./provision <domain> <ip> <services> <docker_registries> <provider>", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := checkArgCount(form); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	provisionTimeout, err := strconv.Atoi(timeoutInSeconds)
	if err != nil {
		fmt.Printf("Error: %s.", err)
		os.Exit(1)
	}

	silentCommandRunner := &runner.ConcreteCmdRunner{
		Stdout:  ioutil.Discard,
		Stderr:  ioutil.Discard,
		Timeout: time.Duration(provisionTimeout) * time.Second,
	}
	p := &provisioner.Provisioner{
		Cert: &cert.Cert{},
		CmdRunner: &runner.ConcreteCmdRunner{
			Stdout:  os.Stdout,
			Stderr:  os.Stderr,
			Timeout: time.Duration(provisionTimeout) * time.Second,
		},
		FS:       &fs.FS{},
		Commands: buildCommands(silentCommandRunner),

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
		case *runner.TimeoutError:
			fmt.Printf("Timed out after %s seconds.\n", timeoutInSeconds)
			os.Exit(1)
		default:
			os.Exit(1)
		}
	}

	w.WriteHeader(http.StatusNoContent)
	fmt.Fprintln(w, "")
}

func stop(w http.ResponseWriter, r *http.Request) {
	cmdRunner := &runner.ConcreteCmdRunner{}
	cmdRunner.Output("/var/vcap/stop")
}

func buildCommands(commandRunner provisioner.CmdRunner) []provisioner.Command {
	providerAgnostic := []provisioner.Command{
		&commands.DisableUAAHSTS{
			WebXMLPath: "/var/vcap/packages/uaa/tomcat/conf/web.xml",
		},
		&commands.ConfigureDnsmasq{
			Domain:     os.Args[1],
			ExternalIP: os.Args[2],
			FS:         &fs.FS{},
			CmdRunner:  commandRunner,
		},
		&commands.ConfigureGardenDNS{
			FS:        &fs.FS{},
			CmdRunner: commandRunner,
		},
		&commands.SetupApi{
			CmdRunner: commandRunner,
			FS:        &fs.FS{},
		},
		&commands.ReplaceDomain{
			CmdRunner: commandRunner,
			FS:        &fs.FS{},
			NewDomain: os.Args[1],
		},
		&commands.SetupCFDot{
			CmdRunner: commandRunner,
			FS:        &fs.FS{},
		},
	}

	const (
		httpPort      = "80"
		httpsPort     = "443"
		sshPort       = "22"
		sshProxyPort  = "2222"
		tcpPortLower  = 61001
		tcpPortHigher = 61100
	)

	forAwsProvider := []provisioner.Command{
		&commands.CloseAllPorts{
			CmdRunner: commandRunner,
		},
		&commands.OpenPort{
			CmdRunner: commandRunner,
			Port:      httpPort,
		},
		&commands.OpenPort{
			CmdRunner: commandRunner,
			Port:      httpsPort,
		},
		&commands.OpenPort{
			CmdRunner: commandRunner,
			Port:      sshPort,
		},
		&commands.OpenPort{
			CmdRunner: commandRunner,
			Port:      sshProxyPort,
		},
	}

	for p := tcpPortLower; p <= tcpPortHigher; p++ {
		forAwsProvider = append(forAwsProvider, &commands.OpenPort{
			CmdRunner: commandRunner,
			Port:      strconv.Itoa(p),
		})
	}

	if isAwsProvisioner() {
		return append(providerAgnostic, forAwsProvider...)
	} else {
		return providerAgnostic
	}
}

func isAwsProvisioner() bool {
	return os.Args[5] == "aws"
}

func handlerStatus(w http.ResponseWriter, r *http.Request) {
	exists, err := fileExists("/run/pcfdev-healthcheck")
	if err != nil {
		fmt.Fprintf(w, fmt.Sprintf(`{"error":{"message":"%s"}}`, err))
	}

	if exists {
		fmt.Fprintf(w, `{"status":"Running"}`)
	} else {
		fmt.Fprintf(w, `{"status":"Unprovisioned"}`)
	}
}

func replaceSecrets(w http.ResponseWriter, r *http.Request) {
	uaaFilePath := "/var/vcap/jobs/uaa/config/uaa.yml"
	uaaCredentialReplacement := &usecases.UaaCredentialReplacement{}

	uaaBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		serverError(w)
		return
	}

	var request struct {
		Password string `json:"password"`
	}

	if err := json.Unmarshal(uaaBytes, &request); err != nil {
		errorHandler(w, "Failed to parse password field from request", http.StatusBadRequest)
		return
	}

	insecureConfig, err := ioutil.ReadFile(uaaFilePath)
	if err != nil {
		serverError(w)
		return
	}

	secureConfig, err := uaaCredentialReplacement.ReplaceUaaConfigAdminCredentials(string(insecureConfig), request.Password)

	if err != nil {
		serverError(w)
		return
	}

	ioutil.WriteFile(uaaFilePath, []byte(secureConfig), 0644)
}

func main() {
	http.HandleFunc("/replace-secrets", replaceSecrets)
	http.HandleFunc("/status", handlerStatus)
	http.HandleFunc("/provision", provision)
	http.HandleFunc("/stop", stop)

	http.ListenAndServe("localhost:8090", nil)
}
