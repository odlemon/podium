package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"podium/internal/api/handlers"
	"podium/internal/models"
)

func (h *Handler) HandleScale(w http.ResponseWriter, r *http.Request) {
	log.Println("Received service scaling request")

	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		handlers.RespondWithError(w, http.StatusBadRequest, "Service ID is required")
		log.Println("Service ID is missing in the request")
		return
	}

	var req models.ServiceScaleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request: %v", err))
		log.Printf("Error decoding request: %v", err)
		return
	}

	log.Printf("Scaling service %s to %d replicas", id, req.Replicas)

	service, err := h.store.GetService(id)
	if err != nil {
		handlers.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("Service not found: %v", err))
		log.Printf("Error retrieving service from database: %v", err)
		return
	}

	currentReplicas := len(service.ContainerIDs)

	if req.Replicas > currentReplicas {
		for i := currentReplicas; i < req.Replicas; i++ {
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
	} else if req.Replicas < currentReplicas {
		containersToRemove := service.ContainerIDs[req.Replicas:currentReplicas]
		service.ContainerIDs = service.ContainerIDs[:req.Replicas]

		for _, containerID := range containersToRemove {
			if err := h.runtime.StopContainer(r.Context(), containerID); err != nil {
				log.Printf("Warning: Failed to stop container %s: %v", containerID, err)
			}

			if err := h.runtime.DeleteContainer(r.Context(), containerID); err != nil {
				log.Printf("Warning: Failed to delete container %s: %v", containerID, err)
			}

			if err := h.store.DeleteContainer(containerID); err != nil {
				log.Printf("Warning: Failed to delete container from database: %v", err)
			}
		}
	}

	service.Replicas = req.Replicas
	service.UpdatedAt = time.Now()

	if err := h.store.UpdateService(service); err != nil {
		handlers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to update service: %v", err))
		log.Printf("Error updating service in database: %v", err)
		return
	}

	handlers.RespondWithJSON(w, http.StatusOK, service)
	log.Printf("Service %s scaled to %d replicas", id, req.Replicas)
}