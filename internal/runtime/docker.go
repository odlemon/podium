package runtime

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"podium/internal/models"
)

type DockerRuntime struct {
	client *client.Client
}

func NewDockerRuntime() (*DockerRuntime, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %w", err)
	}

	return &DockerRuntime{
		client: cli,
	}, nil
}

func (d *DockerRuntime) CreateContainer(ctx context.Context, spec models.Container) error {

	portBindings := nat.PortMap{}
	exposedPorts := nat.PortSet{}

	for _, port := range spec.Ports {
		containerPort := fmt.Sprintf("%d/tcp", port.ContainerPort)
		hostPort := fmt.Sprintf("%d", port.HostPort)

		natPort, err := nat.NewPort("tcp", fmt.Sprintf("%d", port.ContainerPort))
		if err != nil {
			return fmt.Errorf("invalid port mapping: %w", err)
		}

		portBindings[natPort] = []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: hostPort,
			},
		}
		exposedPorts[natPort] = struct{}{}
	}

	env := []string{}
	for k, v := range spec.Env {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	resources := container.Resources{}
	if spec.Resources.CPULimit > 0 {
		resources.NanoCPUs = int64(spec.Resources.CPULimit * 1e9)
	}
	if spec.Resources.MemoryLimit > 0 {
		resources.Memory = spec.Resources.MemoryLimit
	}

	var restartPolicy container.RestartPolicy
	switch spec.RestartPolicy {
	case "Always":
		restartPolicy = container.RestartPolicy{Name: "always"}
	case "OnFailure":
		restartPolicy = container.RestartPolicy{Name: "on-failure"}
	case "Never":
		restartPolicy = container.RestartPolicy{Name: "no"}
	default:
		restartPolicy = container.RestartPolicy{Name: "no"}
	}

	containerConfig := &container.Config{
		Image:        spec.Image,
		Cmd:          spec.Command,
		Env:          env,
		ExposedPorts: exposedPorts,
	}

	hostConfig := &container.HostConfig{
		PortBindings: portBindings,
		Resources:    resources,
		RestartPolicy: restartPolicy,
	}

	resp, err := d.client.ContainerCreate(
		ctx,
		containerConfig,
		hostConfig,
		nil,
		nil,
		spec.Name,
	)
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	if resp.ID != spec.ID {
		return fmt.Errorf("container ID mismatch: expected %s, got %s", spec.ID, resp.ID)
	}

	return nil
}

func (d *DockerRuntime) StartContainer(ctx context.Context, id string) error {
	if err := d.client.ContainerStart(ctx, id, types.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	return nil
}

func (d *DockerRuntime) StopContainer(ctx context.Context, id string) error {
	timeout := 30 * time.Second
	if err := d.client.ContainerStop(ctx, id, &timeout); err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}

	return nil
}

func (d *DockerRuntime) DeleteContainer(ctx context.Context, id string) error {
	if err := d.client.ContainerRemove(ctx, id, types.ContainerRemoveOptions{
		Force: true,
	}); err != nil {
		return fmt.Errorf("failed to remove container: %w", err)
	}

	return nil
}

func (d *DockerRuntime) GetContainerStatus(ctx context.Context, id string) (models.ContainerState, error) {
	resp, err := d.client.ContainerInspect(ctx, id)
	if err != nil {
		return "", fmt.Errorf("failed to inspect container: %w", err)
	}

	if resp.State.Running {
		return models.ContainerStateRunning, nil
	} else if resp.State.ExitCode == 0 {
		return models.ContainerStateSucceeded, nil
	} else {
		return models.ContainerStateFailed, nil
	}
}

func (d *DockerRuntime) GetContainerLogs(ctx context.Context, id string) (string, error) {
	options := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       "100", 
	}

	logs, err := d.client.ContainerLogs(ctx, id, options)
	if err != nil {
		return "", fmt.Errorf("failed to get container logs: %w", err)
	}
	defer logs.Close()

	logBytes, err := io.ReadAll(logs)
	if err != nil {
		return "", fmt.Errorf("failed to read container logs: %w", err)
	}

	return string(logBytes), nil
}