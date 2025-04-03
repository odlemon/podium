package container

import (
	"podium/internal/store"
)

type Handler struct {
	store *store.BoltStore
}

func NewHandler(store *store.BoltStore) *Handler {
	return &Handler{
		store: store,
	}
}