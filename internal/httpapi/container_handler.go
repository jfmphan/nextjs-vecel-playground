package httpapi

import (
	"net/http"

	"homeinventory/internal/service"
)

// ContainerHandler serves the container (storage location) endpoints.
type ContainerHandler struct {
	locations *service.LocationService
}

func NewContainerHandler(locations *service.LocationService) *ContainerHandler {
	return &ContainerHandler{locations: locations}
}

func (h *ContainerHandler) List(w http.ResponseWriter, r *http.Request) {
	containers, err := h.locations.List(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	counts, err := h.locations.ItemCounts(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	out := make([]containerResponse, 0, len(containers))
	for _, c := range containers {
		out = append(out, newContainerResponse(c, counts[c.ID]))
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *ContainerHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req containerRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, err)
		return
	}
	container, err := h.locations.Create(r.Context(), req.toInput())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, newContainerResponse(*container, 0))
}

func (h *ContainerHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := idParam(r)
	if err != nil {
		writeError(w, err)
		return
	}
	container, err := h.locations.Get(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	counts, err := h.locations.ItemCounts(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, newContainerResponse(*container, counts[id]))
}

func (h *ContainerHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := idParam(r)
	if err != nil {
		writeError(w, err)
		return
	}
	var req containerRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, err)
		return
	}
	container, err := h.locations.Update(r.Context(), id, req.toInput())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, newContainerResponse(*container, 0))
}

func (h *ContainerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := idParam(r)
	if err != nil {
		writeError(w, err)
		return
	}
	if err := h.locations.Delete(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
