// Copyright (c) 2024 Christopher Wolf. All rights reserved.
// This software is proprietary and confidential.
// Unauthorized copying of this file, via any medium, is strictly prohibited.

package kafka

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/wolfchristopher/thoth/models"
)

// MockKafkaWriter captures messages in memory for testing.
type MockKafkaWriter struct {
	Messages []kafka.Message
}

func (m *MockKafkaWriter) WriteMessages(ctx context.Context, msgs ...kafka.Message) error {
	m.Messages = append(m.Messages, msgs...)
	return nil
}

func (m *MockKafkaWriter) Close() error {
	return nil
}

// MockGenerateTransaction returns a static transaction with the ID "mock-id".
func MockGenerateTransaction() models.Transaction {
	return models.Transaction{
		ID:        "mock-id",
		Timestamp: "2024-10-27T12:34:56",
		Amount:    99.99,
		Currency:  "USD",
		Customer: models.Customer{
			Name:  "Mock Customer",
			Email: "mock@example.com",
		},
		Items: []models.Item{
			{ProductID: "1234", Quantity: 2, Price: 49.99},
		},
		Status:        "Completed",
		PromotionCode: nil,
		Discount:      nil,
	}
}

func TestKafkaProducer(t *testing.T) {
	mockWriter := &MockKafkaWriter{}

	t.Run("Should send 100 messages to Kafka", func(t *testing.T) {
		StartKafkaProducer(mockWriter, MockGenerateTransaction)

		if len(mockWriter.Messages) != 100 {
			t.Fatalf("Expected 100 messages, but got %d", len(mockWriter.Messages))
		}
	})

	t.Run("Should send correctly formatted XML messages", func(t *testing.T) {
		mockWriter.Messages = nil

		StartKafkaProducer(mockWriter, MockGenerateTransaction)

		for _, msg := range mockWriter.Messages {
			var transaction models.Transaction
			err := xml.Unmarshal(msg.Value, &transaction)
			if err != nil {
				t.Fatalf("Failed to unmarshal XML: %v", err)
			}

			if transaction.ID != "mock-id" {
				t.Fatalf("Expected transaction ID 'mock-id', but got '%s'", transaction.ID)
			}
		}
	})

	t.Run("Should capture all keys correctly", func(t *testing.T) {
		mockWriter.Messages = nil

		StartKafkaProducer(mockWriter, MockGenerateTransaction)

		for i, msg := range mockWriter.Messages {
			expectedKey := []byte(fmt.Sprintf("key-%d", i))
			if !bytes.Equal(msg.Key, expectedKey) {
				t.Fatalf("Expected key '%s', but got '%s'", expectedKey, msg.Key)
			}
		}
	})
}
