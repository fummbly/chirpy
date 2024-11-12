package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, req)

	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, req *http.Request) {

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())))

}

func main() {
	port := "8080"
	directory := "."
	cfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}
	servMux := http.NewServeMux()
	servMux.HandleFunc("/healthz", handlerReadiness)
	servMux.HandleFunc("/metrics", cfg.handlerMetrics)
	servMux.HandleFunc("/reset", cfg.handlerReset)
	appHandler := http.StripPrefix("/app", http.FileServer(http.Dir(directory)))
	servMux.Handle("/app/", cfg.middlewareMetricInc(appHandler))
	server := http.Server{
		Handler: servMux,
		Addr:    ":" + port,
	}

	log.Printf("Serving on port: %s, from directory %s", port, directory)
	log.Fatal(server.ListenAndServe())

}
