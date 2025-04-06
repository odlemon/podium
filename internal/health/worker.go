package health

import (
	"context"
	"log"
	"time"

	"podium/internal/models"
	"podium/internal/runtime"
	"podium/internal/store"
)

type Worker struct {
	store       *store.BoltStore
	runtime     runtime.Runtime
	interval    time.Duration
	maxRestarts int
	stopCh      chan struct{}
}

func NewWorker(store *store.BoltStore, runtime runtime.Runtime, interval time.Duration, maxRestarts int) *Worker {
	return &Worker{
		store:       store,
		runtime:     runtime,
		interval:    interval,
		maxRestarts: maxRestarts,
		stopCh:      make(chan struct{}),
	}
}

func (w *Worker) Start() {
	go func() {
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()

		for {
			select {
			case <-w.stopCh:
				log.Println("Health check worker stopped")
				return
			case <-ticker.C:
				w.checkContainers()
			}
		}
	}()
	log.Println("Health check worker started")
}

func (w *Worker) Stop() {
	close(w.stopCh)
}

func (w *Worker) checkContainers() {
	log.Println("Running health checks on containers")
	
	containers, err := w.store.ListContainers()
	if err != nil {
		log.Printf("Error listing containers for health check: %v", err)
		return
	}
	
	checker := NewChecker()
	
	for _, container := range containers {
		if container.State != models.ContainerStateRunning {
			continue
		}
		
		log.Printf("Checking health of container: %s", container.ID)
		
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		state, err := w.runtime.GetContainerStatus(ctx, container.ID)
		cancel()
		
		if err != nil {
			log.Printf("Error checking container status: %v", err)
			continue
		}
		
		if state != models.ContainerStateRunning {
			log.Printf("Container %s is not running (state: %s)", container.ID, state)
			
			if container.RestartPolicy == "Always" || container.RestartPolicy == "OnFailure" {
				w.restartContainer(container)
			} else {
				log.Printf("Not restarting container %s due to restart policy: %s", container.ID, container.RestartPolicy)
			}
			continue
		}
		
		if container.HealthCheck != nil {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			status, err := checker.Check(ctx, container)
			cancel()
			
			now := time.Now()
			container.Health.LastChecked = now
			
			if err != nil {
				log.Printf("Health check failed for container %s: %v", container.ID, err)
				container.Health.Status = models.HealthStatusUnhealthy
				container.Health.LastFailure = now
				container.Health.FailureCount++
				container.Health.ConsecutiveFail++
				
				if container.Health.ConsecutiveFail >= container.HealthCheck.FailureThreshold {
					log.Printf("Container %s failed health check threshold (%d/%d)", 
						container.ID, container.Health.ConsecutiveFail, container.HealthCheck.FailureThreshold)
					
					if container.RestartPolicy == "Always" || container.RestartPolicy == "OnFailure" {
						w.restartContainer(container)
					}
				}
			} else {
				if status == models.HealthStatusHealthy {
					container.Health.Status = models.HealthStatusHealthy
					container.Health.LastSuccess = now
					container.Health.SuccessCount++
					container.Health.ConsecutiveFail = 0
				} else {
					container.Health.Status = status
				}
			}
			
			if err := w.store.UpdateContainer(container); err != nil {
				log.Printf("Failed to update container health state: %v", err)
			}
		}
	}
	
	log.Println("Container health checks completed")
}

func (w *Worker) restartContainer(container models.Container) {
	log.Printf("Restarting container: %s", container.ID)
	
	if w.maxRestarts > 0 && container.RestartCount >= w.maxRestarts {
		log.Printf("Container %s has exceeded maximum restart count (%d/%d), not restarting", 
			container.ID, container.RestartCount, w.maxRestarts)
		return
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	err := w.runtime.StopContainer(ctx, container.ID)
	cancel()
	
	if err != nil {
		log.Printf("Error stopping container %s: %v", container.ID, err)
	}
	
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	err = w.runtime.StartContainer(ctx, container.ID)
	cancel()
	
	if err != nil {
		log.Printf("Error starting container %s: %v", container.ID, err)
		return
	}
	
	container.State = models.ContainerStateRunning
	container.RestartCount++
	now := time.Now()
	container.StartedAt = &now
	
	if err := w.store.UpdateContainer(container); err != nil {
		log.Printf("Failed to update container state after restart: %v", err)
	}
	
	log.Printf("Container %s restarted successfully (restart count: %d)", container.ID, container.RestartCount)
}