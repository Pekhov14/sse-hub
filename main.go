package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type Client chan string

var (
	clientsMu sync.Mutex
	clients   = make(map[Client]bool)
)

func main() {
	http.HandleFunc("/events", sseHandler)
	http.HandleFunc("/publish", publishHandler)

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func sseHandler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// Headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	client := make(Client)

	// Add client
	clientsMu.Lock()
	clients[client] = true
	clientsMu.Unlock()

	// Remove client on close
	notify := w.(http.CloseNotifier).CloseNotify()
	go func() {
		<-notify
		clientsMu.Lock()
		delete(clients, client)
		clientsMu.Unlock()
		close(client)
	}()

	// Listen for messages
	for msg := range client {
		fmt.Fprintf(w, "data: %s\n\n", msg)
		flusher.Flush()
	}
}

func publishHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Use POST", http.StatusMethodNotAllowed)
		return
	}

	var payload map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Convert to JSON string for broadcasting
	data, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Error encoding data", http.StatusInternalServerError)
		return
	}

	fmt.Println("payload", string(data))

	// Send to all clients
	clientsMu.Lock()
	for client := range clients {
		select {
		case client <- string(data):
		default:
			// if client is slow or closed
			delete(clients, client)
			close(client)
		}
	}
	clientsMu.Unlock()

	w.WriteHeader(http.StatusNoContent)
}
