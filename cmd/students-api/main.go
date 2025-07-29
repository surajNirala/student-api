package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/surajNirala/student-api/internal/config"
)

func main() {
	// Load Config
	cfg := config.MustLoad()
	// Database Setup
	// Setup Router
	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to students api"))
	})
	// Setup Server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	fmt.Printf("Server started %s", cfg.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Failed to start server.")
	}
}
