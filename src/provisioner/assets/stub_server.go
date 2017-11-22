package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Args[1]

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Response from port " + port + " stub server"))
	})

	if err := http.ListenAndServe(":" + port, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
