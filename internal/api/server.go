package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
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
	s.router.HandleFunc("/health", s.handleHealth).Methods("GET")
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) Start(addr string) error {
	return http.ListenAndServe(addr, s.router)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "ok",
		"name":   "Podium",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}