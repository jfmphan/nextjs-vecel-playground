package httpapi

import (
	"net/http"

	"homeinventory/internal/service"
)

// ItemHandler serves the item endpoints.
type ItemHandler struct {
	items *service.InventoryService
}

func NewItemHandler(items *service.InventoryService) *ItemHandler {
	return &ItemHandler{items: items}
}

func (h *ItemHandler) List(w http.ResponseWriter, r *http.Request) {
	filter, err := parseItemFilter(r)
	if err != nil {
		writeError(w, err)
		return
	}
	items, err := h.items.List(r.Context(), filter)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, newItemResponses(items))
}

func (h *ItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req itemRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, err)
		return
	}
	item, err := h.items.Create(r.Context(), req.toInput())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, newItemResponse(*item))
}

func (h *ItemHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := idParam(r)
	if err != nil {
		writeError(w, err)
		return
	}
	item, err := h.items.Get(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, newItemResponse(*item))
}

func (h *ItemHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := idParam(r)
	if err != nil {
		writeError(w, err)
		return
	}
	var req itemRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, err)
		return
	}
	item, err := h.items.Update(r.Context(), id, req.toInput())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, newItemResponse(*item))
}

func (h *ItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := idParam(r)
	if err != nil {
		writeError(w, err)
		return
	}
	if err := h.items.Delete(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
