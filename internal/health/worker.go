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
	store      *store.BoltStore
	runtime    runtime.Runtime
	interval   time.Duration
	maxRetries int
	done       chan struct{}
}

func NewWorker(store *store.BoltStore, runtime runtime.Runtime, interval time.Duration, maxRetries int) *Worker {
	if interval == 0 {
		interval = 30 * time.Second
	}
	
	if maxRetries == 0 {
		maxRetries = 3
	}
	
	return &Worker{
		store:      store,
		runtime:    runtime,
		interval:   interval,
		maxRetries: maxRetries,
		done:       make(chan struct{}),
	}
}

func (w *Worker) Start() {
	log.Println("Starting health check worker with interval:", w.interval)
	
	ticker := time.NewTicker(w.interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				w.checkContainers()
			case <-w.done:
				ticker.Stop()
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	log.Println("Stopping health check worker")
	close(w.done)
}

func (w *Worker) checkContainers() {
	log.Println("Running health checks on containers")
	
	containers, err := w.store.ListContainers()
	if err != nil {
		log.Printf("Error listing containers for health check: %v", err)
		return
	}
	
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
		}
	}
	
	log.Println("Container health checks completed")
}

func (w *Worker) restartContainer(container models.Container) {
	log.Printf("Attempting to restart container: %s", container.ID)
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	err := w.runtime.StartContainer(ctx, container.ID)
	if err != nil {
		log.Printf("Failed to restart container %s: %v", container.ID, err)
		
		container.State = models.ContainerStateFailed
		if err := w.store.UpdateContainer(container); err != nil {
			log.Printf("Failed to update container state: %v", err)
		}
		
		return
	}
	
	log.Printf("Container %s restarted successfully", container.ID)
	
	container.State = models.ContainerStateRunning
	now := time.Now()
	container.StartedAt = &now
	
	if err := w.store.UpdateContainer(container); err != nil {
		log.Printf("Failed to update container state: %v", err)
	}
}