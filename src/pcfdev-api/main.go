package main

import (
	"fmt"
	"net/http"
	"os"
)

const ErrorRouteNotFound = "ROUTE_NOT_FOUND"

func fileExists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, fmt.Sprintf(`{"error":{"message":"%s"}}`, ErrorRouteNotFound))
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

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/status", handlerStatus)
	http.ListenAndServe(":8090", nil)
}