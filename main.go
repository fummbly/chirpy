package main

import (
	"log"
	"net/http"
)

func main() {
	port := "8080"
	directory := "."
	servMux := http.NewServeMux()
	servMux.Handle("/", http.FileServer(http.Dir(directory)))
	server := http.Server{
		Handler: servMux,
		Addr:    ":" + port,
	}

	log.Printf("Serving on port: %s, from directory %s", port, directory)
	log.Fatal(server.ListenAndServe())

}
