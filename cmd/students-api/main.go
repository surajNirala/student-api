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
	"github.com/surajNirala/student-api/internal/storage"
	"github.com/surajNirala/student-api/routes"

	// "gorm.io/driver/mysql"
	"github.com/surajNirala/student-api/internal/storage/mysql"
	"github.com/surajNirala/student-api/internal/storage/sqlite"
)

func main() {
	// Load Config
	cfg := config.MustLoad()
	var storage storage.Storage
	var err error
	// Decide storage based on config
	if cfg.MySQL != nil {
		storage, err = mysql.MysqlConnect(cfg)
		if err != nil {
			log.Fatal("MySQL connection error:", err)
		}
	} else if cfg.StoragePath != "" {
		storage, err = sqlite.New(cfg)
		if err != nil {
			log.Fatal("SQLite connection error:", err)
		}
	} else {
		log.Fatal("No valid database configuration found")
	}
	// Database Setup
	// storage, err := sqlite.New(cfg)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// slog.Info("Storage Intilized", slog.String("env", cfg.Env), slog.String("Version", "1.0.0"))
	// cfg := config.MustLoadMySQL()

	// // Database Setup
	// storage, err := mysql.MysqlConnect(cfg)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	slog.Info("Storage Intilized", slog.String("env", cfg.Env), slog.String("Version", "1.0.0"))

	// Setup Router
	router := http.NewServeMux()
	routes.RouteLoad(router, storage)
	// Setup Server
	server := http.Server{
		Addr:    cfg.HTTPServer.Addr,
		Handler: router,
	}

	//Gressfully shut down server
	slog.Info("Server Started", slog.String("address", cfg.HTTPServer.Addr))
	fmt.Printf("Server Started %s", cfg.HTTPServer.Addr)
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
