package container

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"podium/internal/api/handlers"
)

func (h *Handler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := h.store.GetContainer(id)
	if err != nil {
		handlers.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("Container not found: %v", err))
		return
	}

	if err := h.runtime.DeleteContainer(r.Context(), id); err != nil {
		handlers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to delete container from Docker: %v", err))
		return
	}

	if err := h.store.DeleteContainer(id); err != nil {
		handlers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to delete container from store: %v", err))
		return
	}

	handlers.RespondWithJSON(w, http.StatusNoContent, nil)
}