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
	log.Println("Received container creation request")
	
	var req models.ContainerCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		handlers.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
		return
	}

	log.Printf("Request decoded successfully: name=%s, image=%s", req.Name, req.Image)

	if req.Name == "" || req.Image == "" {
		log.Println("Error: Name and image are required")
		handlers.RespondWithError(w, http.StatusBadRequest, "Name and image are required")
		return
	}

	log.Println("Creating container object")
	container := models.Container{
		ID:            uuid.New().String(),
		Name:          req.Name,
		Image:         req.Image,
		Command:       req.Command,
		Env:           req.Env,
		Ports:         req.Ports,
		Resources:     req.Resources,
		State:         models.ContainerStatePending,
		NodeID:        "local", // For now, we only support a single node
		CreatedAt:     time.Now(),
		RestartPolicy: req.RestartPolicy,
	}
	log.Printf("Container object created with ID: %s", container.ID)

	log.Println("Calling Docker runtime to create container")
	if err := h.runtime.CreateContainer(r.Context(), container); err != nil {
		log.Printf("Error creating container in Docker: %v", err)
		handlers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create container in Docker: %v", err))
		return
	}
	log.Println("Container created successfully in Docker")

	log.Println("Storing container in BoltDB")
	if err := h.store.CreateContainer(container); err != nil {
		log.Printf("Error storing container in BoltDB: %v", err)
		log.Println("Attempting to delete container from Docker due to storage failure")
		
		if deleteErr := h.runtime.DeleteContainer(r.Context(), container.ID); deleteErr != nil {
			log.Printf("Error deleting container from Docker: %v", deleteErr)
		} else {
			log.Println("Container deleted from Docker successfully")
		}
		
		handlers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to store container: %v", err))
		return
	}
	log.Println("Container stored successfully in BoltDB")

	log.Println("Sending successful response")
	handlers.RespondWithJSON(w, http.StatusCreated, container)
	log.Println("Response sent successfully")
}