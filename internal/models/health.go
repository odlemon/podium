package models

import "time"

type HealthCheckType string

const (
	HealthCheckTypeHTTP    HealthCheckType = "http"
	HealthCheckTypeTCP     HealthCheckType = "tcp"
	HealthCheckTypeCommand HealthCheckType = "command"
)

type HealthCheck struct {
	Type            HealthCheckType `json:"type"`
	Endpoint        string          `json:"endpoint,omitempty"`
	Port            int             `json:"port,omitempty"`
	Command         []string        `json:"command,omitempty"`
	InitialDelay    time.Duration   `json:"initialDelay,omitempty"`
	Interval        time.Duration   `json:"interval,omitempty"`
	Timeout         time.Duration   `json:"timeout,omitempty"`
	SuccessThreshold int             `json:"successThreshold,omitempty"`
	FailureThreshold int             `json:"failureThreshold,omitempty"`
}

type HealthStatus string

const (
	HealthStatusUnknown   HealthStatus = "unknown"
	HealthStatusHealthy   HealthStatus = "healthy"
	HealthStatusUnhealthy HealthStatus = "unhealthy"
)

type HealthState struct {
	Status          HealthStatus `json:"status"`
	LastChecked     time.Time    `json:"lastChecked,omitempty"`
	LastSuccess     time.Time    `json:"lastSuccess,omitempty"`
	LastFailure     time.Time    `json:"lastFailure,omitempty"`
	SuccessCount    int          `json:"successCount"`
	FailureCount    int          `json:"failureCount"`
	ConsecutiveFail int          `json:"consecutiveFail"`
}