package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"podium/internal/api/handlers"
	"podium/internal/api/handlers/container"
	"podium/internal/runtime"
	"podium/internal/store"
	"podium/internal/service"
	servicehandler "podium/internal/api/handlers/service"
)

type Server struct {
	router *mux.Router
	store  *store.BoltStore
	runtime runtime.Runtime
	serviceManager service.Manager
}

func NewServer(store *store.BoltStore, runtime runtime.Runtime, serviceManager service.Manager) *Server {
	s := &Server{
		router: mux.NewRouter(),
		store:  store,
		runtime: runtime,
		serviceManager: serviceManager,
	}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.router.HandleFunc("/health", handlers.NewHealthHandler().HandleHealth).Methods("GET")
	
	containerHandler := container.NewHandler(s.store, s.runtime)
	
	s.router.HandleFunc("/api/containers", containerHandler.HandleList).Methods("GET")
	s.router.HandleFunc("/api/containers", containerHandler.HandleCreate).Methods("POST")
	s.router.HandleFunc("/api/containers/{id}", containerHandler.HandleGet).Methods("GET")
	s.router.HandleFunc("/api/containers/{id}", containerHandler.HandleUpdate).Methods("PUT")
	s.router.HandleFunc("/api/containers/{id}", containerHandler.HandleDelete).Methods("DELETE")
	s.router.HandleFunc("/api/containers/{id}/start", containerHandler.HandleStart).Methods("POST")
	s.router.HandleFunc("/api/containers/{id}/stop", containerHandler.HandleStop).Methods("POST")
	s.router.HandleFunc("/api/containers/{id}/logs", containerHandler.HandleLogs).Methods("GET")
	s.router.HandleFunc("/api/containers/{id}/health", containerHandler.HandleHealth).Methods("GET")

	servicehandler.RegisterRoutes(s.router, s.store, s.runtime, s.serviceManager)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) Start(addr string) error {
	return http.ListenAndServe(addr, s.router)
}