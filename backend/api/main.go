package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"garden-planner/internal/config"
	"garden-planner/internal/db"
	"garden-planner/internal/httpapi"
)

func main() {
	cfg := config.FromEnv()
	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	dbConn, err := db.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db connection error: %v", err)
	}
	defer func() {
		_ = dbConn.Close()
	}()

	if cfg.RunMigrations {
		if err := db.Migrate(dbConn, cfg.MigrationsDir); err != nil {
			log.Fatalf("migration failed: %v", err)
		}
		log.Printf("migrations applied from %s", cfg.MigrationsDir)
	}

	companionPath := os.Getenv("COMPANION_DATA_PATH")
	if companionPath == "" {
		companionPath = "datasets/companion_plants_veg.csv"
	}
	companions, err := httpapi.LoadCompanionData(companionPath)
	if err != nil {
		log.Printf("companion data not loaded: %v", err)
	}

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           httpapi.Router(dbConn, companions),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("api listening on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("api server error: %v", err)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("api shutdown error: %v", err)
	}
}
