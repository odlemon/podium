package health

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"podium/internal/models"
)

type Checker struct{}

func NewChecker() *Checker {
	return &Checker{}
}

func (c *Checker) Check(ctx context.Context, container models.Container) (models.HealthStatus, error) {
	if container.HealthCheck == nil {
		return models.HealthStatusHealthy, nil
	}

	switch container.HealthCheck.Type {
	case models.HealthCheckTypeHTTP:
		return c.checkHTTP(ctx, container)
	case models.HealthCheckTypeTCP:
		return c.checkTCP(ctx, container)
	case models.HealthCheckTypeCommand:
		return c.checkCommand(ctx, container)
	default:
		return models.HealthStatusUnknown, fmt.Errorf("unsupported health check type: %s", container.HealthCheck.Type)
	}
}

func (c *Checker) checkHTTP(ctx context.Context, container models.Container) (models.HealthStatus, error) {
	if container.HealthCheck.Endpoint == "" {
		return models.HealthStatusUnknown, fmt.Errorf("HTTP health check requires an endpoint")
	}

	port := container.HealthCheck.Port
	if port == 0 {
		if len(container.Ports) > 0 {
			port = container.Ports[0].ContainerPort
		} else {
			return models.HealthStatusUnknown, fmt.Errorf("HTTP health check requires a port")
		}
	}

	endpoint := container.HealthCheck.Endpoint
	if !strings.HasPrefix(endpoint, "/") {
		endpoint = "/" + endpoint
	}

	u := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("localhost:%d", port),
		Path:   endpoint,
	}


	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return models.HealthStatusUnhealthy, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return models.HealthStatusUnhealthy, fmt.Errorf("HTTP health check failed: %w", err)
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return models.HealthStatusUnhealthy, fmt.Errorf("failed to read HTTP response: %w", err)
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return models.HealthStatusHealthy, nil
	}

	return models.HealthStatusUnhealthy, fmt.Errorf("HTTP health check returned status code: %d", resp.StatusCode)
}

func (c *Checker) checkTCP(ctx context.Context, container models.Container) (models.HealthStatus, error) {
	port := container.HealthCheck.Port
	if port == 0 {
		if len(container.Ports) > 0 {
			port = container.Ports[0].ContainerPort
		} else {
			return models.HealthStatusUnknown, fmt.Errorf("TCP health check requires a port")
		}
	}

	dialer := &net.Dialer{
		Timeout: 5 * time.Second,
	}

	conn, err := dialer.DialContext(ctx, "tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return models.HealthStatusUnhealthy, fmt.Errorf("TCP health check failed: %w", err)
	}
	defer conn.Close()

	return models.HealthStatusHealthy, nil
}

func (c *Checker) checkCommand(ctx context.Context, container models.Container) (models.HealthStatus, error) {
	// For command health checks, we would need to execute a command inside the container
	// This is more complex and would require using the Docker exec API
	// For simplicity, we'll just log that this is not implemented yet
	//Hope you understand
	log.Println("Command health check not implemented yet")
	return models.HealthStatusUnknown, fmt.Errorf("command health check not implemented")
}