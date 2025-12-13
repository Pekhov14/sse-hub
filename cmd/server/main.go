package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"sse/internal/config"
	"sse/internal/sse"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Create SSE server with configured buffer size
	server := sse.NewSSEServer(cfg.ClientBuffer)

	// Setup routes
	http.Handle("/events", server)
	http.HandleFunc("/publish", server.HandlePublish)

	// Create HTTP server
	httpServer := &http.Server{
		Addr: ":" + cfg.Port,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Starting SSE server on port %s", cfg.Port)
		log.Printf("Client buffer size: %d messages", cfg.ClientBuffer)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server gracefully...")

	// Give server 10 seconds to finish current requests
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}
