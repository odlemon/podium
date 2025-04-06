package service

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
	log.Println("Received service creation request")

	var req models.ServiceCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request: %v", err))
		log.Printf("Error decoding request: %v", err)
		return
	}

	log.Printf("Request decoded successfully: name=%s, image=%s, replicas=%d", req.Name, req.Image, req.Replicas)

	service := models.Service{
		ID:            uuid.New().String(),
		Name:          req.Name,
		Image:         req.Image,
		Command:       req.Command,
		Env:           req.Env,
		Ports:         req.Ports,
		Resources:     req.Resources,
		Replicas:      req.Replicas,
		State:         models.ServiceStateCreating,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		RestartPolicy: req.RestartPolicy,
		ContainerIDs:  []string{},
	}

	if err := h.store.CreateService(service); err != nil {
		handlers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create service: %v", err))
		log.Printf("Error storing service in database: %v", err)
		return
	}

	for i := 0; i < service.Replicas; i++ {
		containerName := fmt.Sprintf("%s-%d", service.Name, i)
		container := models.Container{
			ID:            uuid.New().String(),
			Name:          containerName,
			Image:         service.Image,
			Command:       service.Command,
			Ports:         service.Ports,
			Env:           service.Env,
			Resources:     service.Resources,
			State:         models.ContainerStatePending,
			NodeID:        "local",
			CreatedAt:     time.Now(),
			RestartPolicy: service.RestartPolicy,
		}

		if err := h.runtime.CreateContainer(r.Context(), container); err != nil {
			handlers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create container: %v", err))
			log.Printf("Error creating container in Docker: %v", err)
			return
		}

		if err := h.runtime.StartContainer(r.Context(), container.ID); err != nil {
			log.Printf("Warning: Failed to start container %s: %v", container.ID, err)
		} else {
			container.State = models.ContainerStateRunning
			now := time.Now()
			container.StartedAt = &now
		}

		if err := h.store.CreateContainer(container); err != nil {
			log.Printf("Warning: Failed to store container in database: %v", err)
		}

		service.ContainerIDs = append(service.ContainerIDs, container.ID)
	}

	service.State = models.ServiceStateRunning
	service.UpdatedAt = time.Now()

	if err := h.store.UpdateService(service); err != nil {
		log.Printf("Warning: Failed to update service state: %v", err)
	}

	handlers.RespondWithJSON(w, http.StatusCreated, service)
	log.Printf("Service created successfully with ID: %s", service.ID)
}