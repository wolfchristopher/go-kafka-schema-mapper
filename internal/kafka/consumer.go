package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

// StartKafkaConsumer reads messages and maps schema using the provided functions.
func StartKafkaConsumer(
	parseMessageFunc func([]byte) (map[string]interface{}, error),
	mapSchemaFunc func(map[string]interface{}) map[string]string,
) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{"localhost:9092"},
		Topic:          "transactions",
		GroupID:        "transaction-consumer-group",
		MinBytes:       10e3,        // 10KB
		MaxBytes:       10e6,        // 10MB
		CommitInterval: time.Second, // Commit offsets periodically
	})
	defer reader.Close()

	fmt.Println("Kafka consumer started...")

	for {
		message, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatalf("Error reading message from Kafka: %v", err)
		}

		data, err := parseMessageFunc(message.Value)
		if err != nil {
			log.Printf("Failed to parse message: %v", err)
			continue
		}

		schema := mapSchemaFunc(data)

		fmt.Printf("Received Data: %+v\n", data)
		fmt.Printf("Mapped Schema: %+v\n", schema)
	}
}
