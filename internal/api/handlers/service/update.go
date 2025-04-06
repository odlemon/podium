package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"podium/internal/api/handlers"
	"podium/internal/models"
)

func (h *Handler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		handlers.RespondWithError(w, http.StatusBadRequest, "Service ID is required")
		return
	}

	var req models.ServiceCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request: %v", err))
		return
	}

	service, err := h.store.GetService(id)
	if err != nil {
		handlers.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("Service not found: %v", err))
		return
	}

	service.Name = req.Name
	service.Image = req.Image
	service.Command = req.Command
	service.Env = req.Env
	service.Ports = req.Ports
	service.Resources = req.Resources
	service.RestartPolicy = req.RestartPolicy
	service.UpdatedAt = time.Now()

	if err := h.store.UpdateService(service); err != nil {
		handlers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to update service: %v", err))
		return
	}

	handlers.RespondWithJSON(w, http.StatusOK, service)
}