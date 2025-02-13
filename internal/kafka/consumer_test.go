package kafka

import (
	"errors"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
)

func MockParseMessage(value []byte) (map[string]interface{}, error) {
	if string(value) == "error" {
		return nil, errors.New("mock parse error")
	}
	return map[string]interface{}{"message": string(value)}, nil
}

func MockMapSchema() map[string]string {
	return map[string]string{"schema": "mock-schema"}
}

func TestKafkaConsumer(t *testing.T) {
	// Use a channel to simulate message consumption
	messageChan := make(chan kafka.Message)

	go func() {
		for msg := range messageChan {
			data, err := MockParseMessage(msg.Value)
			if err != nil {
				log.Printf("Failed to parse message: %v", err)
				continue
			}
			schema := MockMapSchema()
			fmt.Printf("Received Data: %+v\n", data)
			fmt.Printf("Mapped Schema: %+v\n", schema)
		}
	}()

	t.Run("Should process valid messages", func(t *testing.T) {
		messageChan <- kafka.Message{Value: []byte("valid message")}
		time.Sleep(100 * time.Millisecond) // Allow time for processing
	})

	t.Run("Should handle parsing errors gracefully", func(t *testing.T) {
		messageChan <- kafka.Message{Value: []byte("error")}
		time.Sleep(100 * time.Millisecond)
	})

	t.Run("Should process multiple messages in order", func(t *testing.T) {
		for i := 0; i < 5; i++ {
			messageChan <- kafka.Message{Value: []byte(fmt.Sprintf("message-%d", i))}
		}
		time.Sleep(200 * time.Millisecond)
	})

	// Close the channel after tests are done
	close(messageChan)
}
