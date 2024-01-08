package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gabiSmachado/intents/datamodel"
	"github.com/segmentio/kafka-go"
)

func main() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "test",
		GroupID:  "mygroup",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	var msg datamodel.Intent
	m, err := reader.ReadMessage(context.Background())
	if err != nil {
		log.Fatalf("Error reading message: %v", err)
	}

	err = json.Unmarshal(m.Value, &msg)
	if err != nil {
		log.Fatalf("Error parsing message: %v", err)
	}

	log.Printf("Processed message: %v\n",msg) 
}

