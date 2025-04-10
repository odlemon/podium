package service

import (
	"encoding/json"
	"net/http"

	"podium/internal/models"
	"podium/internal/service"
	"podium/internal/store"
	"podium/internal/utils"
)

type CreateHandler struct {
	Store         store.Store
	ServiceManager service.Manager
}

func (h *CreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var service models.Service
	if err := json.NewDecoder(r.Body).Decode(&service); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if service.Name == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Service name is required")
		return
	}
	if service.Image == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Service image is required")
		return
	}

	service.ID = utils.GenerateID()

	if err := h.Store.CreateService(&service); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create service")
		return
	}

	if err := h.ServiceManager.CreateService(r.Context(), &service); err != nil {
		h.Store.DeleteService(service.ID)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create service containers: "+err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, service)
}