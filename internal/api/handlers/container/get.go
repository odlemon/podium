package container

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"podium/internal/api/handlers"
)

func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	container, err := h.store.GetContainer(id)
	if err != nil {
		handlers.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("Failed to get container: %v", err))
		return
	}

	handlers.RespondWithJSON(w, http.StatusOK, container)
}