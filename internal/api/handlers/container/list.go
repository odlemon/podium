package container

import (
	"fmt"
	"net/http"

	"podium/internal/api/handlers"
	"podium/internal/models"
)

func (h *Handler) HandleList(w http.ResponseWriter, r *http.Request) {
	containers, err := h.store.ListContainers()
	if err != nil {
		handlers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to list containers: %v", err))
		return
	}

	handlers.RespondWithJSON(w, http.StatusOK, map[string][]models.Container{"items": containers})
}