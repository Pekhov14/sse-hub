package keygen

import (
	"fmt"
	"log"

	"sse/internal/auth"
)

// Run generates and prints an API key for the hub.
func Run() {
	key, err := auth.GenerateAPIKey()
	if err != nil {
		log.Fatalf("keygen failed: %v", err)
	}

	fmt.Println("Generated API key:")
	fmt.Println(key)
	fmt.Println()
	fmt.Println("Set it in your environment, for example:")
	fmt.Printf("API_KEY=%s\n", key)
}
