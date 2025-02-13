package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestUpdateKafkaConfig(t *testing.T) {
	validConfig := KafkaConfig{
		Brokers:          "localhost:9092",
		GroupID:          "test-group",
		PreServiceTopic:  "pre_service",
		PostServiceTopic: "post_service",
	}

	t.Run("ValidConfig", func(t *testing.T) {
		jsonData, _ := json.Marshal(validConfig)

		req := httptest.NewRequest(http.MethodPost, "/kafka/config", bytes.NewBuffer(jsonData))
		w := httptest.NewRecorder()

		UpdateKafkaConfig(w, req)

		res := w.Result()
		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %v", res.StatusCode)
		}

		expectedResponse := fmt.Sprintf("Kafka config updated: %+v\n", validConfig)
		if w.Body.String() != expectedResponse {
			t.Errorf("Expected body to be %s, got %s", expectedResponse, w.Body.String())
		}
	})

	t.Run("InvalidConfigFormat", func(t *testing.T) {
		invalidJSON := []byte("invalid json")

		req := httptest.NewRequest(http.MethodPost, "/kafka/config", bytes.NewBuffer(invalidJSON))
		w := httptest.NewRecorder()

		UpdateKafkaConfig(w, req)

		res := w.Result()
		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %v", res.StatusCode)
		}
	})

	t.Run("MethodNotAllowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/kafka/config", nil)
		w := httptest.NewRecorder()

		UpdateKafkaConfig(w, req)

		res := w.Result()
		if res.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("Expected status 405, got %v", res.StatusCode)
		}
	})
}

func TestReceiveSchemaHandler(t *testing.T) {
	t.Run("ValidSchema", func(t *testing.T) {
		validSchema := map[string]interface{}{
			"fields": []interface{}{
				map[string]interface{}{
					"name":     "username",
					"type":     "string",
					"required": true,
				},
				map[string]interface{}{
					"name":     "age",
					"type":     "integer",
					"required": false,
				},
			},
		}

		jsonData, _ := json.Marshal(validSchema)

		req := httptest.NewRequest(http.MethodPost, "/schema", bytes.NewBuffer(jsonData))
		w := httptest.NewRecorder()

		ReceiveSchemaHandler(w, req)

		res := w.Result()
		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %v", res.StatusCode)
		}

		expectedResponse := map[string]string{"message": "Schema received successfully"}
		var actualResponse map[string]string

		if err := json.Unmarshal(w.Body.Bytes(), &actualResponse); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		t.Logf("Expected response: %+v", expectedResponse)
		t.Logf("Actual response: %+v", actualResponse)

		if !reflect.DeepEqual(expectedResponse, actualResponse) {
			t.Errorf("Expected body to be %+v, got %+v", expectedResponse, actualResponse)
		}
	})

	t.Run("InvalidSchemaFormat", func(t *testing.T) {
		invalidJSON := []byte("invalid json")

		req := httptest.NewRequest(http.MethodPost, "/schema", bytes.NewBuffer(invalidJSON))
		w := httptest.NewRecorder()

		ReceiveSchemaHandler(w, req)

		res := w.Result()
		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %v", res.StatusCode)
		}
	})

	t.Run("MethodNotAllowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/schema", nil)
		w := httptest.NewRecorder()

		ReceiveSchemaHandler(w, req)

		res := w.Result()
		if res.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("Expected status 405, got %v", res.StatusCode)
		}
	})
}
