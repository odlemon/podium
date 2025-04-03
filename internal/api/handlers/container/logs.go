package container

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"podium/internal/api/handlers"
)

func (h *Handler) HandleLogs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := h.store.GetContainer(id)
	if err != nil {
		handlers.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("Container not found: %v", err))
		return
	}

	logs, err := h.runtime.GetContainerLogs(r.Context(), id)
	if err != nil {
		handlers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get container logs: %v", err))
		return
	}

	handlers.RespondWithJSON(w, http.StatusOK, map[string]string{"logs": logs})
}