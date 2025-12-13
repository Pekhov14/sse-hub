package sse

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// ServeHTTP handles SSE connections on /events endpoint
func (s *SSEServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Cache-Control")
	w.Header().Set("Access-Control-Allow-Credentials", "false")

	// Create new client channel with configured buffer size
	client := make(Client, s.clientBuffer)
	s.AddClient(client)

	// Handle client disconnect
	ctx := r.Context()
	go func() {
		<-ctx.Done()
		s.RemoveClient(client)
	}()

	// Send messages to client
	for msg := range client {
		_, err := fmt.Fprintf(w, "data: %s\n\n", msg)
		if err != nil {
			break
		}
		flusher.Flush()
	}
}

// HandlePublish handles message publishing on /publish endpoint
func (s *SSEServer) HandlePublish(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight OPTIONS request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Use POST", http.StatusMethodNotAllowed)
		return
	}

	var payload map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	log.Printf("Broadcasting to %d clients: %s", s.GetClientCount(), string(data))
	s.Broadcast(string(data))

	w.WriteHeader(http.StatusNoContent)
}
