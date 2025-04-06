package service

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"podium/internal/api/handlers"
)

func (h *Handler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		handlers.RespondWithError(w, http.StatusBadRequest, "Service ID is required")
		return
	}

	service, err := h.store.GetService(id)
	if err != nil {
		handlers.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("Service not found: %v", err))
		return
	}

	for _, containerID := range service.ContainerIDs {
		if err := h.runtime.StopContainer(r.Context(), containerID); err != nil {
			log.Printf("Warning: Failed to stop container %s: %v", containerID, err)
		}

		if err := h.runtime.DeleteContainer(r.Context(), containerID); err != nil {
			log.Printf("Warning: Failed to delete container %s: %v", containerID, err)
		}

		if err := h.store.DeleteContainer(containerID); err != nil {
			log.Printf("Warning: Failed to delete container from database: %v", err)
		}
	}

	if err := h.store.DeleteService(id); err != nil {
		handlers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to delete service: %v", err))
		return
	}

	handlers.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Service deleted successfully"})
}