package httpapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"homeinventory/internal/auth"
)

// Handlers bundles the constructed handlers and the session manager so the
// router can wire routes without knowing how each handler is built.
type Handlers struct {
	Auth       *AuthHandler
	Items      *ItemHandler
	Containers *ContainerHandler
	Catalog    *CatalogHandler
	Stats      *StatsHandler
	Sessions   *auth.Manager
}

// NewRouter builds the chi router. All routes live under /api/v1: a small public
// group (login/logout/session/health) and a session-protected group for the
// inventory resources.
func NewRouter(h Handlers) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", health)
		r.Post("/login", h.Auth.Login)
		r.Post("/logout", h.Auth.Logout)
		r.Get("/session", h.Auth.Session)

		r.Group(func(r chi.Router) {
			r.Use(requireAuth(h.Sessions))

			r.Route("/items", func(r chi.Router) {
				r.Get("/", h.Items.List)
				r.Post("/", h.Items.Create)
				r.Get("/{id}", h.Items.Get)
				r.Put("/{id}", h.Items.Update)
				r.Delete("/{id}", h.Items.Delete)
			})

			r.Route("/containers", func(r chi.Router) {
				r.Get("/", h.Containers.List)
				r.Post("/", h.Containers.Create)
				r.Get("/{id}", h.Containers.Get)
				r.Put("/{id}", h.Containers.Update)
				r.Delete("/{id}", h.Containers.Delete)
			})

			r.Route("/categories", func(r chi.Router) {
				r.Get("/", h.Catalog.ListCategories)
				r.Post("/", h.Catalog.CreateCategory)
				r.Delete("/{id}", h.Catalog.DeleteCategory)
			})

			r.Get("/tags", h.Catalog.ListTags)
			r.Get("/stats", h.Stats.Get)
		})
	})

	return r
}

func health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
