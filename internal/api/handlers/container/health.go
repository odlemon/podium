package container

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"podium/internal/api/handlers"
)

func (h *Handler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		handlers.RespondWithError(w, http.StatusBadRequest, "Container ID is required")
		return
	}

	container, err := h.store.GetContainer(id)
	if err != nil {
		handlers.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("Container not found: %v", err))
		log.Printf("Error retrieving container from database: %v", err)
		return
	}

	handlers.RespondWithJSON(w, http.StatusOK, container.Health)
}