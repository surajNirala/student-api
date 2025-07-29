package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/surajNirala/student-api/internal/config"
	"github.com/surajNirala/student-api/internal/http/handlers/student"
	"github.com/surajNirala/student-api/internal/storage/sqlite"
)

func main() {
	// Load Config
	cfg := config.MustLoad()
	// Database Setup
	_, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("Storage Intilized", slog.String("env", cfg.Env), slog.String("Version", "1.0.0"))
	// Setup Router
	router := http.NewServeMux()
	router.HandleFunc("POST /api/students/create", student.Create())
	// Setup Server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	//Gressfully shut down server
	slog.Info("Server Started", slog.String("address", cfg.Addr))
	fmt.Printf("Server Started %s", cfg.Addr)
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to start server.")
		}
	}()

	<-done
	slog.Info("Shutting down the server.")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	}
	slog.Info("Server shutdown successfully.")
}
