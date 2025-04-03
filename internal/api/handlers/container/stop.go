package container

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"podium/internal/api/handlers"
	"podium/internal/models"
)

func (h *Handler) HandleStop(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	container, err := h.store.GetContainer(id)
	if err != nil {
		handlers.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("Container not found: %v", err))
		return
	}

	if err := h.runtime.StopContainer(r.Context(), id); err != nil {
		handlers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to stop container: %v", err))
		return
	}

	now := time.Now()
	container.State = models.ContainerStateSucceeded
	container.FinishedAt = &now

	if err := h.store.UpdateContainer(container); err != nil {
		handlers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to update container state: %v", err))
		return
	}

	handlers.RespondWithJSON(w, http.StatusOK, container)
}