package service

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"podium/internal/models"
	"podium/internal/service"
	"podium/internal/store"
)

type HandleCreate struct {
	Store          *store.BoltStore
	ServiceManager service.Manager
}

func (h *HandleCreate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var service models.Service
	if err := json.NewDecoder(r.Body).Decode(&service); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if service.Name == "" {
		respondWithError(w, http.StatusBadRequest, "Service name is required")
		return
	}
	if service.Image == "" {
		respondWithError(w, http.StatusBadRequest, "Service image is required")
		return
	}

	service.ID = generateID()

	if err := h.Store.CreateService(service); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create service")
		return
	}

	if err := h.ServiceManager.CreateService(r.Context(), &service); err != nil {
		h.Store.DeleteService(service.ID)
		respondWithError(w, http.StatusInternalServerError, "Failed to create service containers: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, service)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func generateID() string {
	return "svc-" + strconv.FormatInt(time.Now().UnixNano(), 36)
}