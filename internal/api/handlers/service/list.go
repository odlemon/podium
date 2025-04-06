package service

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"podium/internal/api/handlers"
	"podium/internal/models"
)

func (h *Handler) HandleList(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to list services")
	
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
	
	services, err := h.store.ListServices()
	if err != nil {
		handlers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to list services: %v", err))
		log.Printf("Error listing services from database: %v", err)
		return
	}
	
	var filteredServices []models.Service
	
	if stateFilter != "" {
		log.Printf("Filtering services by state: %s", stateFilter)
		for _, service := range services {
			if string(service.State) == stateFilter {
				filteredServices = append(filteredServices, service)
			}
		}
	} else {
		filteredServices = services
	}
	
	totalCount := len(filteredServices)
	
	if offset >= totalCount {
		offset = 0
	}
	
	endIndex := offset + limit
	if endIndex > totalCount {
		endIndex = totalCount
	}
	
	var paginatedServices []models.Service
	if offset < totalCount {
		paginatedServices = filteredServices[offset:endIndex]
	} else {
		paginatedServices = []models.Service{}
	}
	
	response := map[string]interface{}{
		"items":      paginatedServices,
		"totalCount": totalCount,
		"limit":      limit,
		"offset":     offset,
	}
	
	handlers.RespondWithJSON(w, http.StatusOK, response)
	log.Printf("Listed %d services (filtered from %d total)", len(paginatedServices), totalCount)
}