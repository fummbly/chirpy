package main

import (
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

func main() {
	port := "8080"
	directory := "."
	cfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}
	servMux := http.NewServeMux()
	servMux.HandleFunc("GET /api/healthz", handlerReadiness)
	servMux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	servMux.HandleFunc("POST /admin/reset", cfg.handlerReset)
	appHandler := http.StripPrefix("/app", http.FileServer(http.Dir(directory)))
	servMux.Handle("/app/", cfg.middlewareMetricInc(appHandler))
	server := http.Server{
		Handler: servMux,
		Addr:    ":" + port,
	}

	log.Printf("Serving on port: %s, from directory %s", port, directory)
	log.Fatal(server.ListenAndServe())

}
