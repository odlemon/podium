package container

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"podium/internal/api/handlers"
	"podium/internal/models"
)

func (h *Handler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	log.Println("Container creation request received")

	var req models.ContainerCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		log.Printf("Request decoding failed: %v", err)
		return
	}

	if req.Name == "" || req.Image == "" {
		handlers.RespondWithError(w, http.StatusBadRequest, "Name and image are required")
		log.Println("Validation failed: name or image is missing")
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
		NodeID:        "local",
		CreatedAt:     time.Now(),
		RestartPolicy: req.RestartPolicy,
	}

	if err := h.runtime.CreateContainer(r.Context(), container); err != nil {
		handlers.RespondWithError(w, http.StatusInternalServerError, "Failed to create container in Docker")
		log.Printf("Docker container creation failed: %v", err)
		return
	}

	if err := h.store.CreateContainer(container); err != nil {
		log.Printf("Storing container failed: %v", err)
		if deleteErr := h.runtime.DeleteContainer(r.Context(), container.ID); deleteErr != nil {
			log.Printf("Cleanup failed: could not delete container from Docker: %v", deleteErr)
		}
		handlers.RespondWithError(w, http.StatusInternalServerError, "Failed to store container")
		return
	}

	handlers.RespondWithJSON(w, http.StatusCreated, container)
	log.Printf("Container created successfully: ID=%s", container.ID)
}
