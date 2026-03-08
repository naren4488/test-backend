package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds application configuration loaded from environment.
type Config struct {
	Port   int    // HTTP server port
	DBPath string // Path to SQLite database file
}

const (
	defaultPort   = 8082
	defaultDBPath = "./data/app.db"
)

// Load reads configuration from environment.
// It loads .env from the current directory if present (optional).
// PORT and DB_PATH can override defaults.
func Load() (*Config, error) {
	_ = godotenv.Load() // ignore error if .env is missing

	port := defaultPort
	if v := os.Getenv("PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil && p > 0 {
			port = p
		}
	}

	dbPath := defaultDBPath
	if v := os.Getenv("DB_PATH"); v != "" {
		dbPath = v
	}

	return &Config{
		Port:   port,
		DBPath: dbPath,
	}, nil
}
