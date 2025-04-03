package models

import (
	"time"
)

type ContainerState string

const (
	ContainerStatePending   ContainerState = "pending"
	ContainerStateRunning   ContainerState = "running"
	ContainerStateSucceeded ContainerState = "succeeded"
	ContainerStateFailed    ContainerState = "failed"
)

type PortMapping struct {
	ContainerPort int `json:"containerPort"`
	HostPort      int `json:"hostPort"`
}

type ResourceRequirements struct {
	CPULimit    float64 `json:"cpuLimit"`
	MemoryLimit int64   `json:"memoryLimit"`
}

type Container struct {
	ID            string               `json:"id"`
	Name          string               `json:"name"`
	Image         string               `json:"image"`
	Command       []string             `json:"command,omitempty"`
	Env           map[string]string    `json:"env,omitempty"`
	Ports         []PortMapping        `json:"ports,omitempty"`
	Resources     ResourceRequirements `json:"resources"`
	State         ContainerState       `json:"state"`
	NodeID        string               `json:"nodeId"`
	CreatedAt     time.Time            `json:"createdAt"`
	StartedAt     *time.Time           `json:"startedAt,omitempty"`
	FinishedAt    *time.Time           `json:"finishedAt,omitempty"`
	RestartPolicy string               `json:"restartPolicy"`
}

type ContainerCreateRequest struct {
	Name          string               `json:"name"`
	Image         string               `json:"image"`
	Command       []string             `json:"command,omitempty"`
	Env           map[string]string    `json:"env,omitempty"`
	Ports         []PortMapping        `json:"ports,omitempty"`
	Resources     ResourceRequirements `json:"resources"`
	RestartPolicy string               `json:"restartPolicy"`
}