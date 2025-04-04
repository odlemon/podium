package runtime

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"
	
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"podium/internal/models"
)

type DockerRuntime struct {
	client *client.Client
}

func NewDockerRuntime() (*DockerRuntime, error) {
	log.Println("Initializing Docker client")
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Printf("Error creating Docker client: %v", err)
		return nil, fmt.Errorf("failed to create Docker client: %w", err)
	}
	log.Println("Docker client initialized successfully")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	log.Println("Testing Docker connection")
	_, err = cli.Ping(ctx)
	if err != nil {
		log.Printf("Error connecting to Docker daemon: %v", err)
		return nil, fmt.Errorf("failed to connect to Docker daemon: %w", err)
	}
	log.Println("Docker connection test successful")

	return &DockerRuntime{
		client: cli,
	}, nil
}

func (d *DockerRuntime) CreateContainer(ctx context.Context, spec models.Container) error {
	log.Printf("Creating container: name=%s, image=%s", spec.Name, spec.Image)
	
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	
	log.Println("Setting up port bindings")
	portBindings := nat.PortMap{}
	exposedPorts := nat.PortSet{}

	for _, port := range spec.Ports {
		log.Printf("Setting up port mapping: container=%d, host=%d", port.ContainerPort, port.HostPort)
		natPort, err := nat.NewPort("tcp", fmt.Sprintf("%d", port.ContainerPort))
		if err != nil {
			log.Printf("Error creating port mapping: %v", err)
			return fmt.Errorf("invalid port mapping: %w", err)
		}

		portBindings[natPort] = []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: fmt.Sprintf("%d", port.HostPort),
			},
		}
		exposedPorts[natPort] = struct{}{}
	}

	log.Println("Setting up environment variables")
	env := []string{}
	for k, v := range spec.Env {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	log.Println("Setting up container labels")
	labels := map[string]string{
		"podium.container.id": spec.ID,
	}

	log.Println("Setting up resource limits")
	resources := container.Resources{}
	if spec.Resources.CPULimit > 0 {
		resources.NanoCPUs = int64(spec.Resources.CPULimit * 1e9)
		log.Printf("CPU limit set to: %v cores", spec.Resources.CPULimit)
	}
	if spec.Resources.MemoryLimit > 0 {
		resources.Memory = spec.Resources.MemoryLimit
		log.Printf("Memory limit set to: %v bytes", spec.Resources.MemoryLimit)
	}

	log.Printf("Setting up restart policy: %s", spec.RestartPolicy)
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

	log.Println("Creating container configuration")
	containerConfig := &container.Config{
		Image:        spec.Image,
		Cmd:          spec.Command,
		Env:          env,
		ExposedPorts: exposedPorts,
		Labels:       labels,
	}

	hostConfig := &container.HostConfig{
		PortBindings: portBindings,
		Resources:    resources,
		RestartPolicy: restartPolicy,
	}

	log.Printf("Calling Docker API to create container with ID: %s", spec.ID)
	_, err := d.client.ContainerCreate(
		ctx,
		containerConfig,
		hostConfig,
		nil,
		nil,
		spec.ID,
	)
	if err != nil {
		log.Printf("Error creating container: %v", err)
		return fmt.Errorf("failed to create container: %w", err)
	}

	log.Printf("Container created successfully with ID: %s", spec.ID)
	return nil
}

func (d *DockerRuntime) StartContainer(ctx context.Context, id string) error {
	log.Printf("Starting container: %s", id)
	
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	
	if err := d.client.ContainerStart(ctx, id, container.StartOptions{}); err != nil {
		log.Printf("Error starting container %s: %v", id, err)
		return fmt.Errorf("failed to start container: %w", err)
	}

	log.Printf("Container %s started successfully", id)
	return nil
}

func (d *DockerRuntime) StopContainer(ctx context.Context, id string) error {
	log.Printf("Stopping container: %s", id)
	
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	
	timeoutSeconds := 30
	if err := d.client.ContainerStop(ctx, id, container.StopOptions{Timeout: &timeoutSeconds}); err != nil {
		log.Printf("Error stopping container %s: %v", id, err)
		return fmt.Errorf("failed to stop container: %w", err)
	}

	log.Printf("Container %s stopped successfully", id)
	return nil
}

func (d *DockerRuntime) DeleteContainer(ctx context.Context, id string) error {
	log.Printf("Deleting container: %s", id)
	
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	
	if err := d.client.ContainerRemove(ctx, id, container.RemoveOptions{
		Force: true,
	}); err != nil {
		log.Printf("Error removing container %s: %v", id, err)
		return fmt.Errorf("failed to remove container: %w", err)
	}

	log.Printf("Container %s removed successfully", id)
	return nil
}

func (d *DockerRuntime) GetContainerStatus(ctx context.Context, id string) (models.ContainerState, error) {
	log.Printf("Getting status for container: %s", id)
	
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	
	resp, err := d.client.ContainerInspect(ctx, id)
	if err != nil {
		log.Printf("Error inspecting container %s: %v", id, err)
		return "", fmt.Errorf("failed to inspect container: %w", err)
	}

	var state models.ContainerState
	if resp.State.Running {
		state = models.ContainerStateRunning
	} else if resp.State.ExitCode == 0 {
		state = models.ContainerStateSucceeded
	} else {
		state = models.ContainerStateFailed
	}
	
	log.Printf("Container %s status: %s", id, state)
	return state, nil
}

func (d *DockerRuntime) GetContainerLogs(ctx context.Context, id string) (string, error) {
	log.Printf("Getting logs for container: %s", id)
	
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	
	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       "100",
	}

	logs, err := d.client.ContainerLogs(ctx, id, options)
	if err != nil {
		log.Printf("Error getting logs for container %s: %v", id, err)
		return "", fmt.Errorf("failed to get container logs: %v", err)
	}
	defer logs.Close()

	log.Printf("Reading logs for container %s", id)
	logBytes, err := io.ReadAll(logs)
	if err != nil {
		log.Printf("Error reading logs for container %s: %v", id, err)
		return "", fmt.Errorf("failed to read container logs: %v", err)
	}

	log.Printf("Successfully retrieved logs for container %s (%d bytes)", id, len(logBytes))
	return string(logBytes), nil
}