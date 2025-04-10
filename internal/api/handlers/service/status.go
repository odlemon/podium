package service

import (
	"net/http"

	"github.com/gorilla/mux"
	"podium/internal/service"
	"podium/internal/store"
	"podium/internal/utils"
)

type StatusHandler struct {
	Store         store.Store
	ServiceManager service.Manager
}

func (h *StatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceID := vars["id"]

	if _, err := h.Store.GetService(serviceID); err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Service not found")
		return
	}

	status, err := h.ServiceManager.GetServiceStatus(r.Context(), serviceID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get service status: "+err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, status)
}