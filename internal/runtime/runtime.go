package runtime

import (
	"context"

	"podium/internal/models"
)

type Runtime interface {
	CreateContainer(ctx context.Context, spec models.Container) error
	StartContainer(ctx context.Context, id string) error
	StopContainer(ctx context.Context, id string) error
	DeleteContainer(ctx context.Context, id string) error
	GetContainerStatus(ctx context.Context, id string) (models.ContainerState, error)
	GetContainerLogs(ctx context.Context, id string) (string, error)
}