package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"pcfdev-api/usecases"
)

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
	http.ListenAndServe(":8090", nil)
}
