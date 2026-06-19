package httpapi

import (
	"net/http"

	"homeinventory/internal/auth"
)

// AuthHandler exposes the shared-password login/logout endpoints.
type AuthHandler struct {
	sessions *auth.Manager
}

func NewAuthHandler(sessions *auth.Manager) *AuthHandler {
	return &AuthHandler{sessions: sessions}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Password string `json:"password"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, err)
		return
	}
	if !h.sessions.CheckPassword(req.Password) {
		writeJSON(w, http.StatusUnauthorized, errorBody{Error: "invalid password"})
		return
	}
	h.sessions.Issue(w)
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	h.sessions.Clear(w)
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// Session reports whether the caller is authenticated (used by the frontend to
// decide whether to show the login screen). It is intentionally public.
func (h *AuthHandler) Session(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]bool{"authenticated": h.sessions.Valid(r)})
}
