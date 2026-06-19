// Package app is the composition root: it wires configuration → database →
// repositories → services → handlers → router. Nothing else constructs the
// dependency graph, so the entrypoints (cmd/server, api/index.go) stay tiny.
package app

import (
	"database/sql"
	"fmt"
	"net/http"

	"homeinventory/internal/auth"
	"homeinventory/internal/config"
	"homeinventory/internal/httpapi"
	"homeinventory/internal/service"
	"homeinventory/internal/storage/sqlite"
)

// New builds the fully-wired HTTP handler from configuration. It returns the
// underlying *sql.DB so a long-running server can close it; serverless callers
// can ignore it.
func New(cfg config.Config) (http.Handler, *sql.DB, error) {
	db, err := sqlite.Open(cfg.DatabaseURL, cfg.DatabaseAuthToken)
	if err != nil {
		return nil, nil, err
	}
	if cfg.AutoMigrate {
		if err := sqlite.Migrate(db); err != nil {
			return nil, nil, fmt.Errorf("migrate: %w", err)
		}
	}

	// Storage layer (implements the domain repository interfaces).
	itemRepo := sqlite.NewItemRepo(db)
	containerRepo := sqlite.NewContainerRepo(db)
	categoryRepo := sqlite.NewCategoryRepo(db)
	tagRepo := sqlite.NewTagRepo(db)

	// Application services.
	inventory := service.NewInventoryService(itemRepo, categoryRepo, containerRepo)
	locations := service.NewLocationService(containerRepo)
	catalog := service.NewCatalogService(categoryRepo, tagRepo)
	stats := service.NewStatsService(itemRepo)

	sessions := auth.NewManager(cfg.SessionSecret, cfg.AppPassword, cfg.SecureCookie)

	// HTTP delivery.
	handler := httpapi.NewRouter(httpapi.Handlers{
		Auth:       httpapi.NewAuthHandler(sessions),
		Items:      httpapi.NewItemHandler(inventory),
		Containers: httpapi.NewContainerHandler(locations),
		Catalog:    httpapi.NewCatalogHandler(catalog),
		Stats:      httpapi.NewStatsHandler(stats),
		Sessions:   sessions,
	})
	return handler, db, nil
}
