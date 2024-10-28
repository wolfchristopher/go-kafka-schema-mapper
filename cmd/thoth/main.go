// Copyright (c) 2024 Christopher Wolf. All rights reserved.
// This software is proprietary and confidential.
// Unauthorized copying of this file, via any medium, is strictly prohibited.

package main

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	localkafka "github.com/wolfchristopher/thoth/internal/kafka"
	"github.com/wolfchristopher/thoth/internal/routes"
	"math/rand"
	"net/http"
	"time"
)

// LocalKafkaWriter is a wrapper around Kafka's writer to conform to KafkaWriter interface.
type LocalKafkaWriter struct {
	writer *kafka.Writer
}

func (lw *LocalKafkaWriter) WriteMessages(ctx context.Context, msgs ...kafka.Message) error {
	return lw.writer.WriteMessages(ctx, msgs...)
}

func (lw *LocalKafkaWriter) Close() error {
	return lw.writer.Close()
}

func main() {
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/kafka_config", routes.UpdateKafkaConfig)
	http.HandleFunc("/schema", routes.ReceiveSchemaHandler)

	writer := &LocalKafkaWriter{
		writer: &kafka.Writer{
			Addr:     kafka.TCP("localhost:9092"),
			Topic:    "transactions",
			Balancer: &kafka.LeastBytes{},
		},
	}
	defer func(writer *LocalKafkaWriter) {
		err := writer.Close()
		if err != nil {

		}
	}(writer)

	go localkafka.StartKafkaProducer(writer, localkafka.GenerateTransaction)
	go localkafka.StartKafkaConsumer(
		func(value []byte) (map[string]interface{}, error) {
			return map[string]interface{}{"message": string(value)}, nil
		},
		func(data map[string]interface{}) map[string]string {
			return map[string]string{"schema": "mock-schema"}
		},
	)

	fmt.Println("Starting server on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
