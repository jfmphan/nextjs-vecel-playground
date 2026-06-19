// Command migrate applies all pending database migrations and exits. Use it to
// migrate a remote (Turso) database, which is not auto-migrated on startup:
//
//	DATABASE_URL=libsql://... DATABASE_AUTH_TOKEN=... go run ./cmd/migrate
package main

import (
	"log"

	"homeinventory/internal/config"
	"homeinventory/internal/storage/sqlite"
)

func main() {
	cfg := config.Load()
	db, err := sqlite.Open(cfg.DatabaseURL, cfg.DatabaseAuthToken)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	if err := sqlite.Migrate(db); err != nil {
		log.Fatalf("migrate: %v", err)
	}
	log.Println("migrations applied")
}
