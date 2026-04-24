package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/iict-sust/inventory/internal/services"
)

type StatsHandler struct {
	statsSvc *services.StatsService
}

func NewStatsHandler(statsSvc *services.StatsService) *StatsHandler {
	return &StatsHandler{statsSvc: statsSvc}
}

func (h *StatsHandler) GetPublicStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats, err := h.statsSvc.GetPublicStats(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (h *StatsHandler) GetDashboardStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats, err := h.statsSvc.GetDashboardStats(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
