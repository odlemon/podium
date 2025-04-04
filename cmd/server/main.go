package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"podium/internal/api"
	"podium/internal/runtime"
	"podium/internal/store"
)

func main() {
	fmt.Println("It's Podium baby")

	boltStore, err := store.NewBoltStore("podium.db")
	if err != nil {
		log.Fatalf("Failed to create store: %v", err)
	}
	defer boltStore.Close()

	fmt.Println("BoltDB store initialized successfully")

	dockerRuntime, err := runtime.NewDockerRuntime()
	if err != nil {
		log.Fatalf("Failed to create Docker runtime: %v", err)
	}
	fmt.Println("Docker runtime initialized successfully")

	server := api.NewServer(boltStore, dockerRuntime)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	errChan := make(chan error)
	go func() {
		fmt.Println("API server listening on :8080")
		errChan <- server.Start(":8080")
	}()

	select {
	case err := <-errChan:
		log.Fatalf("Server error: %v", err)
	case sig := <-sigChan:
		fmt.Printf("Received signal: %v, shutting down...\n", sig)
	}

	fmt.Println("Podium stopped")
}