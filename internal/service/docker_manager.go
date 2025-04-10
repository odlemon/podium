package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	
	"podium/internal/models"
	"podium/internal/store"
)

type DockerServiceManager struct {
	dockerClient *client.Client
	store        store.Store
}

func NewDockerServiceManager(dockerClient *client.Client, store store.Store) *DockerServiceManager {
	return &DockerServiceManager{
		dockerClient: dockerClient,
		store:        store,
	}
}

func (m *DockerServiceManager) CreateService(ctx context.Context, service *models.Service) error {
	if service.Replicas <= 0 {
		service.Replicas = 1
	}

	for i := 0; i < service.Replicas; i++ {
		if err := m.createServiceContainer(ctx, service, i); err != nil {
			return fmt.Errorf("failed to create container %d for service %s: %w", i, service.ID, err)
		}
	}

	return nil
}

func (m *DockerServiceManager) createServiceContainer(ctx context.Context, service *models.Service, index int) error {
	containerName := fmt.Sprintf("%s-%d", service.Name, index)

	portBindings := nat.PortMap{}
	exposedPorts := nat.PortSet{}

	for _, port := range service.Ports {
		containerPort := nat.Port(fmt.Sprintf("%d/tcp", port.Container))
		exposedPorts[containerPort] = struct{}{}
		
		hostPort := port.Host
		if index > 0 {
			hostPort = 0
		}
		
		portBindings[containerPort] = []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: strconv.Itoa(hostPort),
			},
		}
	}

	env := []string{}
	for k, v := range service.Environment {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	env = append(env, 
		fmt.Sprintf("PODIUM_SERVICE_ID=%s", service.ID),
		fmt.Sprintf("PODIUM_SERVICE_NAME=%s", service.Name),
		fmt.Sprintf("PODIUM_REPLICA_INDEX=%d", index),
	)

	config := &container.Config{
		Image:        service.Image,
		ExposedPorts: exposedPorts,
		Env:          env,
		Labels: map[string]string{
			"podium.service.id":   service.ID,
			"podium.service.name": service.Name,
			"podium.replica.index": strconv.Itoa(index),
			"podium.managed":      "true",
		},
	}

	if service.HealthCheck != nil {
		interval, _ := time.ParseDuration(service.HealthCheck.Interval)
		timeout, _ := time.ParseDuration(service.HealthCheck.Timeout)

		var healthCmd string
		switch service.HealthCheck.Type {
		case "http":
			healthCmd = fmt.Sprintf("curl -f http://localhost:%d%s || exit 1", 
				service.HealthCheck.Port, 
				service.HealthCheck.Endpoint)
		case "tcp":
			healthCmd = fmt.Sprintf("nc -z localhost %d || exit 1", 
				service.HealthCheck.Port)
		case "command":
			healthCmd = service.HealthCheck.Command
		default:
			healthCmd = fmt.Sprintf("curl -f http://localhost:%d/ || exit 1", 
				service.Ports[0].Container)
		}

		config.Healthcheck = &container.HealthConfig{
			Test:     []string{"CMD-SHELL", healthCmd},
			Interval: interval,
			Timeout:  timeout,
			Retries:  service.HealthCheck.Retries,
		}
	}

	hostConfig := &container.HostConfig{
		PortBindings: portBindings,
		RestartPolicy: container.RestartPolicy{
			Name: convertRestartPolicy(service.RestartPolicy),
		},
	}

	resp, err := m.dockerClient.ContainerCreate(
		ctx,
		config,
		hostConfig,
		nil,
		nil,
		containerName,
	)
	if err != nil {
		return err
	}

	if err := m.dockerClient.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	return nil
}

func convertRestartPolicy(policy *models.RestartPolicy) string {
	if policy == nil {
		return "no"
	}

	switch policy.Type {
	case "always":
		return "always"
	case "on-failure":
		return "on-failure"
	case "unless-stopped":
		return "unless-stopped"
	default:
		return "no"
	}
}

func (m *DockerServiceManager) UpdateService(ctx context.Context, service *models.Service) error {
	containers, err := m.getServiceContainers(ctx, service.ID)
	if err != nil {
		return err
	}
	
	for _, container := range containers {
		if err := m.dockerClient.ContainerStop(ctx, container.ID, container.StopOptions); err != nil {
			return err
		}
		if err := m.dockerClient.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{}); err != nil {
			return err
		}
	}

	return m.CreateService(ctx, service)
}

func (m *DockerServiceManager) DeleteService(ctx context.Context, serviceID string) error {
	containers, err := m.getServiceContainers(ctx, serviceID)
	if err != nil {
		return err
	}

	for _, container := range containers {
		if err := m.dockerClient.ContainerStop(ctx, container.ID, container.StopOptions); err != nil {
			return err
		}
		if err := m.dockerClient.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{}); err != nil {
			return err
		}
	}

	return nil
}

func (m *DockerServiceManager) ScaleService(ctx context.Context, serviceID string, replicas int) error {
	service, err := m.store.GetService(serviceID)
	if err != nil {
		return err
	}

	containers, err := m.getServiceContainers(ctx, serviceID)
	if err != nil {
		return err
	}

	currentCount := len(containers)

	if replicas > currentCount {
		service.Replicas = replicas
		for i := currentCount; i < replicas; i++ {
			if err := m.createServiceContainer(ctx, service, i); err != nil {
				return err
			}
		}
	}

	if replicas < currentCount {
		for i := currentCount - 1; i >= replicas; i-- {
			for _, container := range containers {
				indexStr := container.Labels["podium.replica.index"]
				index, _ := strconv.Atoi(indexStr)
				if index == i {
					if err := m.dockerClient.ContainerStop(ctx, container.ID, container.StopOptions); err != nil {
						return err
					}
					if err := m.dockerClient.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{}); err != nil {
						return err
					}
					break
				}
			}
		}
		
		service.Replicas = replicas
		if err := m.store.UpdateService(service); err != nil {
			return err
		}
	}

	return nil
}

func (m *DockerServiceManager) GetServiceStatus(ctx context.Context, serviceID string) (*ServiceStatus, error) {
	containers, err := m.getServiceContainers(ctx, serviceID)
	if err != nil {
		return nil, err
	}

	service, err := m.store.GetService(serviceID)
	if err != nil {
		return nil, err
	}

	healthyCount := 0
	containerStatuses := make([]ContainerStatus, 0, len(containers))

	for _, container := range containers {
		healthState := "unknown"
		if container.State != nil && container.State.Health != nil {
			healthState = container.State.Health.Status
		}

		if healthState == "healthy" {
			healthyCount++
		}

		containerStatuses = append(containerStatuses, ContainerStatus{
			ID:          container.ID,
			Name:        container.Names[0],
			Status:      container.State.Status,
			HealthState: healthState,
			CreatedAt:   container.Created.Format(time.RFC3339),
			StartedAt:   container.State.StartedAt.Format(time.RFC3339),
		})
	}

	return &ServiceStatus{
		ServiceID:       serviceID,
		DesiredReplicas: service.Replicas,
		CurrentReplicas: len(containers),
		HealthyReplicas: healthyCount,
		Containers:      containerStatuses,
	}, nil
}

func (m *DockerServiceManager) ReconcileServices(ctx context.Context) error {
	services, err := m.store.ListServices()
	if err != nil {
		return err
	}

	for _, service := range services {
		status, err := m.GetServiceStatus(ctx, service.ID)
		if err != nil {
			continue
		}

		if status.CurrentReplicas < service.Replicas {
			if err := m.ScaleService(ctx, service.ID, service.Replicas); err != nil {
				continue
			}
		}

		if status.CurrentReplicas > service.Replicas {
			if err := m.ScaleService(ctx, service.ID, service.Replicas); err != nil {
				continue
			}
		}

		for _, containerStatus := range status.Containers {
			if containerStatus.HealthState == "unhealthy" {
				if err := m.dockerClient.ContainerRestart(ctx, containerStatus.ID, container.StopOptions{}); err != nil {
					continue
				}
			}
		}
	}

	return nil
}

func (m *DockerServiceManager) getServiceContainers(ctx context.Context, serviceID string) ([]types.Container, error) {
	filterArgs := filters.NewArgs()
	filterArgs.Add("label", fmt.Sprintf("podium.service.id=%s", serviceID))

	return m.dockerClient.ContainerList(ctx, types.ContainerListOptions{
		Filters: filterArgs,
		All:     true,
	})
}