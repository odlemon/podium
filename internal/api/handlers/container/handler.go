package container

import (
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