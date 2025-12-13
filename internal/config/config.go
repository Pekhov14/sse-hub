package config

import (
	"os"
	"strconv"
)

// Config holds application configuration
type Config struct {
	Port         string
	ClientBuffer int
	CORSOrigin   string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	cfg := &Config{
		Port:         getEnv("PORT", "8080"),
		ClientBuffer: getEnvInt("CLIENT_BUFFER", 10),
		CORSOrigin:   getEnv("CORS_ORIGIN", "*"),
	}
	return cfg
}

// getEnv gets environment variable with default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt gets environment variable as int with default value
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
