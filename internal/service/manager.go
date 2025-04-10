package service

import (
	"context"

	"podium/internal/models"
)

type Manager interface {
	CreateService(ctx context.Context, service *models.Service) error
	UpdateService(ctx context.Context, service *models.Service) error
	DeleteService(ctx context.Context, serviceID string) error
	ScaleService(ctx context.Context, serviceID string, replicas int) error
	GetServiceStatus(ctx context.Context, serviceID string) (*ServiceStatus, error)
	ReconcileServices(ctx context.Context) error
}

type ServiceStatus struct {
	ServiceID      string
	DesiredReplicas int
	CurrentReplicas int
	HealthyReplicas int
	Containers     []ContainerStatus
}

type ContainerStatus struct {
	ID          string
	Name        string
	Status      string
	HealthState string
	CreatedAt   string
	StartedAt   string
}