package sse

import (
	"sync"
)

// Client represents a connected SSE client
type Client chan string

// SSEServer manages SSE connections and message broadcasting
type SSEServer struct {
	clients      map[Client]bool
	clientsMu    sync.Mutex
	clientBuffer int
}

// NewSSEServer creates a new SSE server instance
func NewSSEServer(clientBuffer int) *SSEServer {
	if clientBuffer <= 0 {
		clientBuffer = 10 // default buffer size
	}
	return &SSEServer{
		clients:      make(map[Client]bool),
		clientBuffer: clientBuffer,
	}
}

// AddClient adds a new client to the server
func (s *SSEServer) AddClient(c Client) {
	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()
	s.clients[c] = true
}

// RemoveClient removes a client from the server and closes its channel
func (s *SSEServer) RemoveClient(c Client) {
	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()
	delete(s.clients, c)
	close(c)
}

// Broadcast sends a message to all connected clients
func (s *SSEServer) Broadcast(msg string) {
	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()
	for client := range s.clients {
		select {
		case client <- msg:
		default:
			// Client channel is full or closed, remove it
			delete(s.clients, client)
			close(client)
		}
	}
}

// GetClientCount returns the number of connected clients
func (s *SSEServer) GetClientCount() int {
	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()
	return len(s.clients)
}
