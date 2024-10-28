// Copyright (c) 2024 Christopher Wolf. All rights reserved.
// This software is proprietary and confidential.
// Unauthorized copying of this file, via any medium, is strictly prohibited.

package main

import (
	"fmt"
	"github.com/segmentio/kafka-go"
	localkafka "github.com/wolfchristopher/thoth/internal/kafka"
	"github.com/wolfchristopher/thoth/internal/routes"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/kafka_config", routes.UpdateKafkaConfig)
	http.HandleFunc("/schema", routes.ReceiveSchemaHandler)

	writer := &localkafka.LocalKafkaWriter{
		Writer: &kafka.Writer{
			Addr:     kafka.TCP("localhost:9092"),
			Topic:    "transactions",
			Balancer: &kafka.LeastBytes{},
		},
	}
	defer func(writer localkafka.KafkaWriter) {
		if err := writer.Close(); err != nil {
			fmt.Printf("Error closing writer: %v\n", err)
		}
	}(writer)

	go localkafka.StartKafkaProducer(writer, localkafka.GenerateTransaction)

	// Start Kafka consumer with mock logic
	go localkafka.StartKafkaConsumer(
		func(value []byte) (map[string]interface{}, error) {
			return map[string]interface{}{"message": string(value)}, nil
		},
		func(data map[string]interface{}) map[string]string {
			return map[string]string{"schema": "mock-schema"}
		},
	)

	fmt.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
