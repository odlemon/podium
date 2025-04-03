package container

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"podium/internal/api/handlers"
	"podium/internal/models"
)

func (h *Handler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	container, err := h.store.GetContainer(id)
	if err != nil {
		handlers.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("Container not found: %v", err))
		return
	}

	var req models.Container
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
		return
	}

	container.State = req.State
	if req.StartedAt != nil {
		container.StartedAt = req.StartedAt
	}
	if req.FinishedAt != nil {
		container.FinishedAt = req.FinishedAt
	}

	if err := h.store.UpdateContainer(container); err != nil {
		handlers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to update container: %v", err))
		return
	}

	handlers.RespondWithJSON(w, http.StatusOK, container)
}