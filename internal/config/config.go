// Package config loads runtime configuration from environment variables.
package config

import (
	"os"
	"strings"
)

// Config holds all runtime configuration. See .env.example for the full list
// and the local-development defaults applied by Load.
type Config struct {
	DatabaseURL       string // "libsql://..." (Turso) or "file:..." (local dev)
	DatabaseAuthToken string // Turso auth token; unused for local file DBs
	SessionSecret     string // HMAC key used to sign session cookies
	AppPassword       string // the single shared password that gates the app
	AutoMigrate       bool   // apply pending migrations on startup
	SecureCookie      bool   // mark the session cookie Secure (HTTPS only)
}

// Load reads configuration from the environment, applying defaults that make
// local development work with zero setup (a WAL SQLite file under ./data).
func Load() Config {
	databaseURL := getEnv("DATABASE_URL",
		"file:./data/inventory.db?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)")
	isLocalFile := strings.HasPrefix(databaseURL, "file:")

	return Config{
		DatabaseURL:       databaseURL,
		DatabaseAuthToken: os.Getenv("DATABASE_AUTH_TOKEN"),
		SessionSecret:     getEnv("SESSION_SECRET", "dev-insecure-secret-change-me"),
		AppPassword:       getEnv("APP_PASSWORD", "changeme"),
		// Local file DBs auto-migrate for zero-setup dev. A remote DB (Turso) is
		// migrated explicitly via `go run ./cmd/migrate` unless AUTO_MIGRATE=1,
		// so production cold starts don't pay for a migration check.
		AutoMigrate:  isLocalFile || os.Getenv("AUTO_MIGRATE") == "1",
		SecureCookie: !isLocalFile, // any non-local deployment is assumed to be HTTPS
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
