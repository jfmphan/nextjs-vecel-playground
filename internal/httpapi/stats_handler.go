package httpapi

import (
	"net/http"
	"strconv"

	"homeinventory/internal/service"
)

// defaultExpiringWindowDays is the look-ahead window for "expiring soon" when
// the client does not specify one.
const defaultExpiringWindowDays = 30

// StatsHandler serves the dashboard summary endpoint.
type StatsHandler struct {
	stats *service.StatsService
}

func NewStatsHandler(stats *service.StatsService) *StatsHandler {
	return &StatsHandler{stats: stats}
}

func (h *StatsHandler) Get(w http.ResponseWriter, r *http.Request) {
	days := defaultExpiringWindowDays
	if v := r.URL.Query().Get("expiringWithinDays"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			days = n
		}
	}
	stats, err := h.stats.Compute(r.Context(), days)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, newStatsResponse(stats))
}
