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

	if err := h.store.DeleteContainer(id); err != nil {
		handlers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to delete container: %v", err))
		return
	}

	handlers.RespondWithJSON(w, http.StatusNoContent, nil)
}