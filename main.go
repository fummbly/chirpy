package main

import (
	"log"
	"net/http"
)

func handlerReadiness(w http.ResponseWriter, req *http.Request) {

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))

}

func main() {
	port := "8080"
	directory := "."
	servMux := http.NewServeMux()
	servMux.HandleFunc("/healthz", handlerReadiness)
	servMux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(directory))))
	server := http.Server{
		Handler: servMux,
		Addr:    ":" + port,
	}

	log.Printf("Serving on port: %s, from directory %s", port, directory)
	log.Fatal(server.ListenAndServe())

}
