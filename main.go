package main

import (
	"fmt"
	"net/http"
)

func main() {
	servMux := http.NewServeMux()
	server := http.Server{
		Handler: servMux,
		Addr:    ":8080",
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Failed to start server: %v", err)
		return
	}

}
