// Package handler is the Vercel serverless entrypoint. Vercel routes every
// request under /api/* to Handler (see vercel.json, which rewrites /api/v1/* to
// this function). The application is built once per warm instance and reused
// across requests, so the database connection is not reopened on every call.
package handler

import (
	"net/http"
	"sync"

	"homeinventory/bootstrap"
)

var (
	buildOnce sync.Once
	router    http.Handler
	buildErr  error
)

// Handler is the exported function Vercel invokes for each request.
func Handler(w http.ResponseWriter, r *http.Request) {
	buildOnce.Do(func() {
		router, buildErr = bootstrap.NewRouter()
	})
	if buildErr != nil {
		http.Error(w, "service unavailable", http.StatusServiceUnavailable)
		return
	}
	router.ServeHTTP(w, r)
}
