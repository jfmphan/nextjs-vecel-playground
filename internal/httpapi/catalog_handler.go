package httpapi

import (
	"net/http"

	"homeinventory/internal/service"
)

// CatalogHandler serves the category and tag endpoints.
type CatalogHandler struct {
	catalog *service.CatalogService
}

func NewCatalogHandler(catalog *service.CatalogService) *CatalogHandler {
	return &CatalogHandler{catalog: catalog}
}

func (h *CatalogHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.catalog.ListCategories(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, newCategoryResponses(categories))
}

func (h *CatalogHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req categoryRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, err)
		return
	}
	category, err := h.catalog.CreateCategory(r.Context(), req.Name)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, namedResponse{ID: category.ID, Name: category.Name})
}

func (h *CatalogHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, err := idParam(r)
	if err != nil {
		writeError(w, err)
		return
	}
	if err := h.catalog.DeleteCategory(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *CatalogHandler) ListTags(w http.ResponseWriter, r *http.Request) {
	tags, err := h.catalog.ListTags(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, newTagResponses(tags))
}
