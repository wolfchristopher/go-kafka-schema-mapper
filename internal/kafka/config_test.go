// Copyright (c) 2024 Christopher Wolf. All rights reserved.
// This software is proprietary and confidential.
// Unauthorized copying of this file, via any medium, is strictly prohibited.

package kafka

import (
	"reflect"
	"testing"
)

func TestMapSchema(t *testing.T) {
	// Behavior: It should map a flat structure with primitive types
	t.Run("Flat structure with primitive types", func(t *testing.T) {
		input := map[string]interface{}{
			"username": "john_doe",
			"age":      30,
			"isActive": true,
		}

		expected := map[string]interface{}{
			"username": "string",
			"age":      "int",
			"isActive": "bool",
		}

		result := MapSchema(input)

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %+v, but got %+v", expected, result)
		}
	})

	// Behavior: It should map nested structures (maps within maps)
	t.Run("Nested structures", func(t *testing.T) {
		input := map[string]interface{}{
			"user": map[string]interface{}{
				"username": "john_doe",
				"profile": map[string]interface{}{
					"age":      30,
					"isActive": true,
				},
			},
		}

		expected := map[string]interface{}{
			"user": map[string]interface{}{
				"username": "string",
				"profile": map[string]interface{}{
					"age":      "int",
					"isActive": "bool",
				},
			},
		}

		result := MapSchema(input)

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %+v, but got %+v", expected, result)
		}
	})

	// Behavior: It should return an empty map when given an empty input
	t.Run("Empty input map", func(t *testing.T) {
		input := map[string]interface{}{}
		expected := map[string]interface{}{}

		result := MapSchema(input)

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %+v, but got %+v", expected, result)
		}
	})

	// Behavior: It should handle mixed types
	t.Run("Mixed types", func(t *testing.T) {
		input := map[string]interface{}{
			"name":     "John Doe",
			"age":      25,
			"isActive": true,
			"meta": map[string]interface{}{
				"tags": []string{"go", "bdd", "testing"},
			},
		}

		expected := map[string]interface{}{
			"name":     "string",
			"age":      "int",
			"isActive": "bool",
			"meta": map[string]interface{}{
				"tags": "[]string",
			},
		}

		result := MapSchema(input)

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %+v, but got %+v", expected, result)
		}
	})

	// Behavior: It should handle nil values gracefully
	t.Run("Nil values", func(t *testing.T) {
		input := map[string]interface{}{
			"username": nil,
		}

		expected := map[string]interface{}{
			"username": "<nil>",
		}

		result := MapSchema(input)

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %+v, but got %+v", expected, result)
		}
	})
}
