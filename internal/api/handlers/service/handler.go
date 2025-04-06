package service

import (
	"github.com/gorilla/mux"
	"podium/internal/runtime"
	"podium/internal/store"
)

type Handler struct {
	store   *store.BoltStore
	runtime runtime.Runtime
}

func NewHandler(store *store.BoltStore, runtime runtime.Runtime) *Handler {
	return &Handler{
		store:   store,
		runtime: runtime,
	}
}

func RegisterRoutes(router *mux.Router, store *store.BoltStore, runtime runtime.Runtime) {
	h := NewHandler(store, runtime)
	
	router.HandleFunc("/api/services", h.HandleList).Methods("GET")
	router.HandleFunc("/api/services", h.HandleCreate).Methods("POST")
	router.HandleFunc("/api/services/{id}", h.HandleGet).Methods("GET")
	router.HandleFunc("/api/services/{id}", h.HandleUpdate).Methods("PUT")
	router.HandleFunc("/api/services/{id}", h.HandleDelete).Methods("DELETE")
	router.HandleFunc("/api/services/{id}/scale", h.HandleScale).Methods("POST")
}