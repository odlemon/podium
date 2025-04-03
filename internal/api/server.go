package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"podium/internal/api/handlers"
	"podium/internal/api/handlers/container"
	"podium/internal/store"
)

type Server struct {
	router *mux.Router
	store  *store.BoltStore
}

func NewServer(store *store.BoltStore) *Server {
	s := &Server{
		router: mux.NewRouter(),
		store:  store,
	}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {

	s.router.HandleFunc("/health", handlers.NewHealthHandler().HandleHealth).Methods("GET")
	
	containerHandler := container.NewHandler(s.store)
	
	s.router.HandleFunc("/api/containers", containerHandler.HandleList).Methods("GET")
	s.router.HandleFunc("/api/containers", containerHandler.HandleCreate).Methods("POST")
	s.router.HandleFunc("/api/containers/{id}", containerHandler.HandleGet).Methods("GET")
	s.router.HandleFunc("/api/containers/{id}", containerHandler.HandleUpdate).Methods("PUT")
	s.router.HandleFunc("/api/containers/{id}", containerHandler.HandleDelete).Methods("DELETE")
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) Start(addr string) error {
	return http.ListenAndServe(addr, s.router)
}