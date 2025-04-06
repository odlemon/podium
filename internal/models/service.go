package models

import "time"


type ServiceState string

const (

	ServiceStateCreating ServiceState = "creating"
	
	ServiceStateRunning ServiceState = "running"

	ServiceStateStopped ServiceState = "stopped"

	ServiceStateFailed ServiceState = "failed"
)

type Service struct {
	ID            string               `json:"id"`
	Name          string               `json:"name"`
	Image         string               `json:"image"`
	Command       []string             `json:"command,omitempty"`
	Env           map[string]string    `json:"env,omitempty"`
	Ports         []PortMapping        `json:"ports,omitempty"`
	Resources     ResourceRequirements `json:"resources"`
	Replicas      int                  `json:"replicas"`
	State         ServiceState         `json:"state"`
	CreatedAt     time.Time            `json:"createdAt"`
	UpdatedAt     time.Time            `json:"updatedAt"`
	RestartPolicy string               `json:"restartPolicy"`
	ContainerIDs  []string             `json:"containerIds,omitempty"`
}

type ServiceCreateRequest struct {
	Name          string               `json:"name"`
	Image         string               `json:"image"`
	Command       []string             `json:"command,omitempty"`
	Env           map[string]string    `json:"env,omitempty"`
	Ports         []PortMapping        `json:"ports,omitempty"`
	Resources     ResourceRequirements `json:"resources"`
	Replicas      int                  `json:"replicas"`
	RestartPolicy string               `json:"restartPolicy"`
}

type ServiceScaleRequest struct {
	Replicas int `json:"replicas"`
}