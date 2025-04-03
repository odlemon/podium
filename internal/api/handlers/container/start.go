package container

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"podium/internal/api/handlers"
	"podium/internal/models"
)

func (h *Handler) HandleStart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	container, err := h.store.GetContainer(id)
	if err != nil {
		handlers.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("Container not found: %v", err))
		return
	}

	if err := h.runtime.StartContainer(r.Context(), id); err != nil {
		handlers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to start container: %v", err))
		return
	}

	now := time.Now()
	container.State = models.ContainerStateRunning
	container.StartedAt = &now

	if err := h.store.UpdateContainer(container); err != nil {
		handlers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to update container state: %v", err))
		return
	}

	handlers.RespondWithJSON(w, http.StatusOK, container)
}