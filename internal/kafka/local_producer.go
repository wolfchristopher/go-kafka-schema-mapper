// Copyright (c) 2024 Christopher Wolf. All rights reserved.
// This software is proprietary and confidential.
// Unauthorized copying of this file, via any medium, is strictly prohibited.

package kafka

import (
	"context"
	"encoding/xml"
	"fmt"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/wolfchristopher/thoth/models"
	"log"
	"math/rand"
	"time"
)

// LocalKafkaWriter KafkaWriter defines the interface to mock Kafka writing.
type LocalKafkaWriter interface {
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

// StartKafkaProducer generates and sends transactions to Kafka.
func StartKafkaProducer(writer LocalKafkaWriter, generateTransaction func() models.Transaction) {
	defer func() {
		if err := writer.Close(); err != nil {
			panic(err)
		}
	}()

	for i := 0; i < 100; i++ {
		// Use only the injected function to generate transactions
		transaction := generateTransaction()

		xmlData, err := xml.MarshalIndent(transaction, "  ", "    ")
		if err != nil {
			log.Fatalf("Error marshaling XML: %v", err)
		}

		err = writer.WriteMessages(context.Background(), kafka.Message{
			Key:   []byte(fmt.Sprintf("key-%d", i)),
			Value: xmlData,
		})
		if err != nil {
			log.Fatalf("Failed to write message to Kafka: %v", err)
		}

		fmt.Printf("Sent transaction %d to Kafka\n", i+1)
	}
}

// GenerateTransaction Generate a transaction with random data.
func GenerateTransaction() models.Transaction {
	return models.Transaction{
		ID:        uuid.New().String(),
		Timestamp: generateRandomTimestamp(),
		Amount:    float64(rand.Intn(990)+10) + rand.Float64(),
		Currency:  randomChoice([]string{"USD", "EUR", "JPY"}),
		Customer: models.Customer{
			Name:  randomChoice([]string{"Alice", "Bob", "Charlie", "David", "Eve"}),
			Email: randomChoice([]string{"alice@example.com", "bob@example.com", "charlie@example.com"}),
		},
		Items: generateItems(),
		Status: randomChoice([]string{
			"Completed", "Pending", "Cancelled",
		}),
		PromotionCode: randomPromotionCode(),
		Discount:      randomDiscount(),
	}
}

func generateRandomTimestamp() string {
	now := time.Now()
	randomDuration := time.Duration(rand.Intn(30*24)) * time.Hour
	return now.Add(-randomDuration).Format("2006-01-02T15:04:05")
}

func generateItems() []models.Item {
	numItems := rand.Intn(3) + 1
	items := make([]models.Item, numItems)
	for i := 0; i < numItems; i++ {
		items[i] = models.Item{
			ProductID: fmt.Sprintf("%d", rand.Intn(9000)+1000),
			Quantity:  rand.Intn(10) + 1,
			Price:     float64(rand.Intn(500)) + rand.Float64(),
		}
	}
	return items
}

func randomPromotionCode() *string {
	if rand.Float32() < 0.3 {
		code := fmt.Sprintf("PROMO-%d", rand.Intn(100))
		return &code
	}
	return nil
}

func randomDiscount() *float64 {
	if rand.Float32() < 0.2 {
		discount := float64(rand.Intn(30)) + rand.Float64()
		return &discount
	}
	return nil
}

func randomChoice(choices []string) string {
	return choices[rand.Intn(len(choices))]
}
