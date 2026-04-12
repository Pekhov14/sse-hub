package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// sse - Server-sent events
// sk  - Secret key
const apiKeyPrefix = "sse_sk_"

// GenerateAPIKey creates a random API key for publisher authentication.
func GenerateAPIKey() (string, error) {
	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		return "", fmt.Errorf("generate api key: %w", err)
	}

	return apiKeyPrefix + hex.EncodeToString(secret), nil
}
