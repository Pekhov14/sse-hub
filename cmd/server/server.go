package server

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

// Run starts the SSE hub server.
func Run() {
	cfg := config.LoadConfig()

	server := sse.NewSSEServer(cfg.ClientBuffer)

	http.Handle("/events", server)
	http.HandleFunc("/publish", server.HandlePublish)

	httpServer := &http.Server{
		Addr: ":" + cfg.Port,
	}

	go func() {
		log.Printf("Starting SSE server on port %s", cfg.Port)
		log.Printf("Client buffer size: %d messages", cfg.ClientBuffer)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}
