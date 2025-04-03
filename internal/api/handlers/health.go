package handlers

import (
	"encoding/json"
	"net/http"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "ok",
		"name":   "Podium",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}