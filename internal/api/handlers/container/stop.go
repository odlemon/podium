package container

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"podium/internal/api/handlers"
	"podium/internal/models"
)

func (h *Handler) HandleStop(w http.ResponseWriter, r *http.Request) {
	log.Println("Received container stop request")
	
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		handlers.RespondWithError(w, http.StatusBadRequest, "Container ID is required")
		log.Println("Container ID is missing in the request")
		return
	}
	
	log.Printf("Stopping container with ID: %s", id)
	
	container, err := h.store.GetContainer(id)
	if err != nil {
		handlers.RespondWithError(w, http.StatusInternalServerError, "Failed to get container")
		log.Printf("Error retrieving container from database: %v", err)
		return
	}
	
	if container.ID == "" {
		handlers.RespondWithError(w, http.StatusNotFound, "Container not found")
		log.Printf("Container with ID %s not found", id)
		return
	}

	err = h.runtime.StopContainer(r.Context(), id)
	if err != nil {
		handlers.RespondWithError(w, http.StatusInternalServerError, "Failed to stop container")
		log.Printf("Error stopping container in Docker: %v", err)
		return
	}
	
	container.State = models.ContainerStateSucceeded
	now := time.Now()
	container.FinishedAt = &now
	
	err = h.store.UpdateContainer(container)
	if err != nil {
		log.Printf("Warning: Failed to update container state in database: %v", err)
	}
	
	handlers.RespondWithJSON(w, http.StatusOK, container)
	log.Printf("Container %s stopped successfully", id)
}