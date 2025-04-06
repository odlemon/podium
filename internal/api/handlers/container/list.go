package container

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"podium/internal/api/handlers"
	"podium/internal/models"
)

func (h *Handler) HandleList(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to list containers")
	
	query := r.URL.Query()
	stateFilter := query.Get("state")
	
	limitStr := query.Get("limit")
	offsetStr := query.Get("offset")
	
	var limit, offset int
	var err error
	
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			limit = 100 
		}
	} else {
		limit = 100 
	}
	
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			offset = 0 
		}
	}
	
	containers, err := h.store.ListContainers()
	if err != nil {
		handlers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to list containers: %v", err))
		log.Printf("Error listing containers from database: %v", err)
		return
	}
	
	var filteredContainers []models.Container
	
	if stateFilter != "" {
		log.Printf("Filtering containers by state: %s", stateFilter)
		for _, container := range containers {
			if string(container.State) == stateFilter {
				filteredContainers = append(filteredContainers, container)
			}
		}
	} else {
		filteredContainers = containers
	}
	
	totalCount := len(filteredContainers)
	
	if offset >= totalCount {
		offset = 0
	}
	
	endIndex := offset + limit
	if endIndex > totalCount {
		endIndex = totalCount
	}
	
	var paginatedContainers []models.Container
	if offset < totalCount {
		paginatedContainers = filteredContainers[offset:endIndex]
	} else {
		paginatedContainers = []models.Container{}
	}
	
	response := map[string]interface{}{
		"items":      paginatedContainers,
		"totalCount": totalCount,
		"limit":      limit,
		"offset":     offset,
	}
	
	handlers.RespondWithJSON(w, http.StatusOK, response)
	log.Printf("Listed %d containers (filtered from %d total)", len(paginatedContainers), totalCount)
}