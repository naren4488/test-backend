package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds application configuration loaded from environment.
type Config struct {
	Port             int      // HTTP server port
	DBPath           string   // Path to SQLite database file
	JWTSecret        string   // Secret for signing JWTs
	JWTExpiryHours   int      // Token validity in hours
	CORSAllowedOrigins []string // Allowed CORS origins (e.g. http://localhost:3000)
}

const (
	defaultPort         = 8082
	defaultDBPath       = "./data/app.db"
	defaultJWTExpiryHours = 24
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

	jwtSecret := os.Getenv("JWT_SECRET")
	jwtExpiry := defaultJWTExpiryHours
	if v := os.Getenv("JWT_EXPIRY_HOURS"); v != "" {
		if e, err := strconv.Atoi(v); err == nil && e > 0 {
			jwtExpiry = e
		}
	}

	corsOrigins := []string{
		"http://localhost:3000",
		"http://localhost:5173",
		"http://127.0.0.1:3000",
		"http://127.0.0.1:5173",
	}
	if v := os.Getenv("CORS_ORIGINS"); v != "" {
		corsOrigins = strings.Split(strings.TrimSpace(v), ",")
		for i := range corsOrigins {
			corsOrigins[i] = strings.TrimSpace(corsOrigins[i])
		}
	}

	return &Config{
		Port:               port,
		DBPath:             dbPath,
		JWTSecret:          jwtSecret,
		JWTExpiryHours:     jwtExpiry,
		CORSAllowedOrigins: corsOrigins,
	}, nil
}
