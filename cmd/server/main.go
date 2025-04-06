package main

import (
	"log"
	"time"
	
	"podium/internal/api"
	"podium/internal/health"
	"podium/internal/runtime"
	"podium/internal/store"
)

func main() {
	log.Println("It's Podium baby")
	
	boltStore, err := store.NewBoltStore("podium.db")
	if err != nil {
		log.Fatalf("Failed to create store: %v", err)
	}
	defer boltStore.Close()
	
	dockerRuntime, err := runtime.NewDockerRuntime()
	if err != nil {
		log.Fatalf("Failed to create Docker runtime: %v", err)
	}
	
	server := api.NewServer(boltStore, dockerRuntime)
	
	healthWorker := health.NewWorker(boltStore, dockerRuntime, 30*time.Second, 3)
	healthWorker.Start()
	defer healthWorker.Stop()
	
	log.Println("Starting server on :8080")
	if err := server.Start(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}