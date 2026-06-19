// Package bootstrap wires the application together and exposes it as a plain
// http.Handler. It lives outside internal/ on purpose: the Vercel serverless
// entrypoint (api/index.go) is compiled by Vercel outside the module tree, which
// makes Go forbid it from importing internal/* directly. The entrypoint imports
// this package instead, and bootstrap (a normal in-module package) is free to
// reach into internal/.
package bootstrap

import (
	"net/http"

	"homeinventory/internal/app"
	"homeinventory/internal/config"
)

// NewRouter builds the application HTTP handler from the loaded configuration.
// The underlying *sql.DB is owned by the handler for the lifetime of the
// process, so it is intentionally not returned here.
func NewRouter() (http.Handler, error) {
	router, _, err := app.New(config.Load())
	return router, err
}
