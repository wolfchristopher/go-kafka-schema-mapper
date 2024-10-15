package main

import (
	"fmt"
	"github.com/wolfchristopher/thoth/internal/routes"
	"net/http"
)

func main() {
	http.HandleFunc("/kafka_config", routes.UpdateKafkaConfig)
	http.HandleFunc("/schema", routes.ReceiveSchemaHandler)

	fmt.Println("Starting server on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
