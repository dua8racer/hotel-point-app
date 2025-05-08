package config

import (
	"os"
	"strconv"
)

type Config struct {
	Server struct {
		Port string
	}
	MongoDB struct {
		URI      string
		Database string
	}
	JWT struct {
		Secret      string
		ExpiryHours int
	}
}

func NewConfig() *Config {
	cfg := &Config{}

	// Server configuration
	cfg.Server.Port = getEnv("SERVER_PORT", "8080")

	// MongoDB configuration
	cfg.MongoDB.URI = getEnv("MONGO_URI", "mongodb://myuser_xxx:mypassword_xxx@192.168.1.29:31039/condotel?authSource=admin")
	cfg.MongoDB.Database = getEnv("MONGO_DB", "condotel")

	// JWT configuration
	cfg.JWT.Secret = getEnv("JWT_SECRET", "Fr3eP@le5t1n3!!!")
	cfg.JWT.ExpiryHours, _ = strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))

	return cfg
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
