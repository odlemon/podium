package service

import (
	"net/http"

	"github.com/gorilla/mux"
	"podium/internal/service"
	"podium/internal/store"
)

type HandleStatus struct {
	Store          *store.BoltStore
	ServiceManager service.Manager
}

func (h *HandleStatus) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceID := vars["id"]

	if _, err := h.Store.GetService(serviceID); err != nil {
		respondWithError(w, http.StatusNotFound, "Service not found")
		return
	}

	status, err := h.ServiceManager.GetServiceStatus(r.Context(), serviceID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get service status: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, status)
}