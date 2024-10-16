package kafka

import (
	"reflect"
)

func MapSchema(data map[string]interface{}) map[string]interface{} {
	schema := make(map[string]interface{})
	for key, value := range data {
		// Check for nil values before using reflect.TypeOf
		if value == nil {
			schema[key] = "<nil>"
			continue
		}
		// Store the type of each field
		schema[key] = reflect.TypeOf(value).String()

		// If the value is a nested structure, recursively discover its schema
		if reflect.TypeOf(value).Kind() == reflect.Map {
			schema[key] = MapSchema(value.(map[string]interface{}))
		}
	}
	return schema
}
