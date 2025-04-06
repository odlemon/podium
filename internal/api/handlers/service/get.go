package service

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"podium/internal/api/handlers"
)

func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		handlers.RespondWithError(w, http.StatusBadRequest, "Service ID is required")
		return
	}

	service, err := h.store.GetService(id)
	if err != nil {
		handlers.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("Service not found: %v", err))
		log.Printf("Error retrieving service from database: %v", err)
		return
	}

	handlers.RespondWithJSON(w, http.StatusOK, service)
}