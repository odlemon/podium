package service

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"podium/internal/service"
	"podium/internal/store"
	"podium/internal/utils"
)

type ScaleRequest struct {
	Replicas int `json:"replicas"`
}

type ScaleHandler struct {
	Store         store.Store
	ServiceManager service.Manager
}

func (h *ScaleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceID := vars["id"]

	var req ScaleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Replicas < 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "Replicas count must be non-negative")
		return
	}

	service, err := h.Store.GetService(serviceID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Service not found")
		return
	}

	if err := h.ServiceManager.ScaleService(r.Context(), serviceID, req.Replicas); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to scale service: "+err.Error())
		return
	}

	service.Replicas = req.Replicas
	if err := h.Store.UpdateService(service); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update service")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, service)
}