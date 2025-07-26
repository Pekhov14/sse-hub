package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type Client chan string

type SSEServer struct {
	clients   map[Client]bool
	clientsMu sync.Mutex
}

func NewSSEServer() *SSEServer {
	return &SSEServer{
		clients: make(map[Client]bool),
	}
}

func (s *SSEServer) AddClient(c Client) {
	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()
	s.clients[c] = true
}

func (s *SSEServer) RemoveClient(c Client) {
	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()
	delete(s.clients, c)
	close(c)
}

func (s *SSEServer) Broadcast(msg string) {
	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()
	for client := range s.clients {
		select {
		case client <- msg:
		default:
			delete(s.clients, client)
			close(client)
		}
	}
}

func (s *SSEServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	client := make(Client, 10)
	s.AddClient(client)

	ctx := r.Context()
	go func() {
		<-ctx.Done()
		s.RemoveClient(client)
	}()

	for msg := range client {
		_, err := fmt.Fprintf(w, "data: %s\n\n", msg)
		if err != nil {
			break
		}
		flusher.Flush()
	}
}

func (s *SSEServer) HandlePublish(w http.ResponseWriter, r *http.Request) {
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

	log.Println("Broadcasting payload:", string(data))
	s.Broadcast(string(data))

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	server := NewSSEServer()

	http.Handle("/events", server)
	http.HandleFunc("/publish", server.HandlePublish)

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
