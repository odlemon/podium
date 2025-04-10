package service

import (
	"github.com/gorilla/mux"
	"podium/internal/runtime"
	"podium/internal/service"
	"podium/internal/store"
)

type Handler struct {
	store          *store.BoltStore
	runtime        runtime.Runtime
	serviceManager service.Manager
}

func NewHandler(store *store.BoltStore, runtime runtime.Runtime, serviceManager service.Manager) *Handler {
	return &Handler{
		store:          store,
		runtime:        runtime,
		serviceManager: serviceManager,
	}
}

func RegisterRoutes(router *mux.Router, store *store.BoltStore, runtime runtime.Runtime, serviceManager service.Manager) {
	h := NewHandler(store, runtime, serviceManager)
	
	router.HandleFunc("/api/services", h.HandleList).Methods("GET")
	router.HandleFunc("/api/services", h.HandleCreate).Methods("POST")
	router.HandleFunc("/api/services/{id}", h.HandleGet).Methods("GET")
	router.HandleFunc("/api/services/{id}", h.HandleUpdate).Methods("PUT")
	router.HandleFunc("/api/services/{id}", h.HandleDelete).Methods("DELETE")
	router.HandleFunc("/api/services/{id}/scale", h.HandlerScale).Methods("POST")
	router.HandleFunc("/api/services/{id}/status", h.HandleStatus).Methods("GET")
}