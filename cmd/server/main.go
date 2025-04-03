package main

import (
	"fmt"
	"log"

	"github.com/odlemon/podium/internal/store"
)

func main() {
	fmt.Println("Starting Podium - Container Orchestration Tool")

	boltStore, err := store.NewBoltStore("podium.db")
	if err != nil {
		log.Fatalf("Failed to create store: %v", err)
	}
	defer boltStore.Close()

	fmt.Println("BoltDB store initialized successfully")
}