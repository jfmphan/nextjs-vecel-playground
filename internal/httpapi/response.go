// Package httpapi is the HTTP delivery layer: it maps requests to service calls,
// translates between JSON DTOs and domain/service types, and maps errors to
// status codes. It contains no business logic.
package httpapi

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"homeinventory/internal/domain"
)

type errorBody struct {
	Error string `json:"error"`
}

// writeJSON encodes v as JSON with the given status code.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("httpapi: encode response: %v", err)
	}
}

// writeError maps an error to an HTTP status using the domain sentinels. Unknown
// errors are treated as internal and their detail is hidden from the client.
func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		writeJSON(w, http.StatusNotFound, errorBody{Error: err.Error()})
	case errors.Is(err, domain.ErrValidation):
		writeJSON(w, http.StatusBadRequest, errorBody{Error: err.Error()})
	case errors.Is(err, domain.ErrConflict):
		writeJSON(w, http.StatusConflict, errorBody{Error: err.Error()})
	default:
		log.Printf("httpapi: internal error: %v", err)
		writeJSON(w, http.StatusInternalServerError, errorBody{Error: "internal server error"})
	}
}

// decodeJSON reads and decodes a JSON request body into dst, returning a
// validation error on malformed input.
func decodeJSON(r *http.Request, dst any) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return domain.Invalid("invalid JSON body")
	}
	return nil
}
