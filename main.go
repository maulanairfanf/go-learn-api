package main

import (
	"context"
	"log"
	"log/slog"
	"myapi/db"
	"myapi/routes"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	db.Init()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := routes.InitializeRoutes()

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Jalanin server di goroutine terpisah
	go func() {
		slog.Info("Starting server", "port", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Tunggu sinyal berhenti (Ctrl+C / SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")

	// Kasih waktu 5 detik buat request yang lagi jalan selesai
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	db.Close()
	slog.Info("Server exited")
}
