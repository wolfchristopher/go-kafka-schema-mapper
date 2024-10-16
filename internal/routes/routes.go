package routes

import (
	"encoding/json"
	"fmt"
	"github.com/wolfchristopher/thoth/internal/kafka"
	"log"
	"net/http"
)

type KafkaConfig struct {
	Brokers          string `json:"brokers"`
	GroupID          string `json:"group_id"`
	PreServiceTopic  string `json:"pre_service_topic"`
	PostServiceTopic string `json:"post_service_topic"`
	SecurityProtocol string `json:"security_protocol,omitempty"`
	SASLMechanism    string `json:"sasl_mechanism,omitempty"`
	SASLUsername     string `json:"sasl_username,omitempty"`
	SASLPassword     string `json:"sasl_password,omitempty"`
	SSLCaLocation    string `json:"ssl_ca_location,omitempty"`
	AutoOffsetReset  string `json:"auto_offset_reset,omitempty"`
}

var currentConfig KafkaConfig

func UpdateKafkaConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&currentConfig)
		if err != nil {
			http.Error(w, "Invalid config format", http.StatusBadRequest)
			return
		}
		_, err = fmt.Fprintf(w, "Kafka config updated: %+v\n", currentConfig)
		if err != nil {
			return
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func ReceiveSchemaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var schema map[string]interface{}

	err := json.NewDecoder(r.Body).Decode(&schema)
	if err != nil {
		log.Printf("Failed to decode JSON response: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Received schema: %+v\n", schema)

	schema = kafka.MapSchema(schema)

	log.Printf("Schema Mapped: %+v\n", schema)

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]string{"message": "Schema received successfully"})
	if err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
		http.Error(w, "Failed to send the response", http.StatusInternalServerError)
		return
	}
}
