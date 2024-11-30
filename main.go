package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/fummbly/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	database       *database.Queries
	platform       string
	secret         string
}

func (cfg *apiConfig) middlewareMetricInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, req)

	})
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("Error getting sql server: %v\n", err)
		return
	}

	port := "8080"
	directory := "."
	cfg := apiConfig{
		fileserverHits: atomic.Int32{},
		database:       database.New(db),
		platform:       os.Getenv("PLATFORM"),
		secret:         os.Getenv("SECRET"),
	}

	servMux := http.NewServeMux()
	servMux.HandleFunc("GET /api/healthz", handlerReadiness)
	servMux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	servMux.HandleFunc("GET /api/chirps", cfg.handleGetChirps)
	servMux.HandleFunc("GET /api/chirps/{chirp_id}", cfg.handleGetChirp)
	servMux.HandleFunc("POST /admin/reset", cfg.handlerReset)
	servMux.HandleFunc("POST /api/chirps", cfg.handleAddChirp)
	servMux.HandleFunc("POST /api/users", cfg.handleAddUser)
	servMux.HandleFunc("POST /api/login", cfg.handleLogin)
	appHandler := http.StripPrefix("/app", http.FileServer(http.Dir(directory)))
	servMux.Handle("/app/", cfg.middlewareMetricInc(appHandler))
	server := http.Server{
		Handler: servMux,
		Addr:    ":" + port,
	}

	log.Printf("Serving on port: %s, from directory %s", port, directory)
	log.Fatal(server.ListenAndServe())

}
