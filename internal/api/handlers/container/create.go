package container

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"podium/internal/api/handlers"
	"podium/internal/models"
)

func (h *Handler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var req models.ContainerCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
		return
	}

	if req.Name == "" || req.Image == "" {
		handlers.RespondWithError(w, http.StatusBadRequest, "Name and image are required")
		return
	}

	container := models.Container{
		ID:            uuid.New().String(),
		Name:          req.Name,
		Image:         req.Image,
		Command:       req.Command,
		Env:           req.Env,
		Ports:         req.Ports,
		Resources:     req.Resources,
		State:         models.ContainerStatePending,
		NodeID:        "local", // only supports  single node
		CreatedAt:     time.Now(),
		RestartPolicy: req.RestartPolicy,
	}

	if err := h.store.CreateContainer(container); err != nil {
		handlers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create container: %v", err))
		return
	}

	handlers.RespondWithJSON(w, http.StatusCreated, container)
}