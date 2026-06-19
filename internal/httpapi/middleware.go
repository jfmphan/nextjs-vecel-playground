package httpapi

import (
	"net/http"

	"homeinventory/internal/auth"
)

// requireAuth rejects requests that do not carry a valid session cookie.
func requireAuth(sessions *auth.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !sessions.Valid(r) {
				writeJSON(w, http.StatusUnauthorized, errorBody{Error: "unauthorized"})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
