package main

import (
	"log"

	"github.com/axe-junction/axe-server/vtc-service/internal/handlers"
)

func main() {
	if err := handlers.StartKafkaEstimateWorker(); err != nil {
		log.Fatalf("Failed to start Kafka estimate worker: %v", err)
	}
}