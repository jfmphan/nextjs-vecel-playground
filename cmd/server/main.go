// Command server runs the inventory API as a standalone HTTP server for local
// development. In production the same application is served by the Vercel
// function in api/index.go.
package main

import (
	"log"
	"net/http"
	"os"

	"homeinventory/internal/app"
	"homeinventory/internal/config"
)

func main() {
	handler, db, err := app.New(config.Load())
	if err != nil {
		log.Fatalf("startup: %v", err)
	}
	defer db.Close()

	addr := ":" + envOr("PORT", "8080")
	log.Printf("home inventory API listening on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
