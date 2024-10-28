// Copyright (c) 2024 Christopher Wolf. All rights reserved.
// This software is proprietary and confidential.
// Unauthorized copying of this file, via any medium, is strictly prohibited.

package kafka

import (
	"bytes"
	"reflect"
	"testing"
)

func TestMapSchema(t *testing.T) {
	t.Run("Given a flat structure, it should map schema correctly", func(t *testing.T) {
		input := map[string]interface{}{
			"id":     "123",
			"age":    30,
			"active": true,
		}
		expected := map[string]interface{}{
			"id":     "string",
			"age":    "int",
			"active": "bool",
		}
		result := MapSchema(input)

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Given a nested structure, it should recursively map the schema", func(t *testing.T) {
		input := map[string]interface{}{
			"user": map[string]interface{}{
				"name": "Alice",
				"info": map[string]interface{}{
					"age": 25,
				},
			},
		}
		expected := map[string]interface{}{
			"user": map[string]interface{}{
				"name": "string",
				"info": map[string]interface{}{
					"age": "int",
				},
			},
		}
		result := MapSchema(input)

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})
}

func TestJSONToMap(t *testing.T) {
	t.Run("Given valid JSON, it should parse into a map", func(t *testing.T) {
		data := []byte(`{"name": "Alice", "age": 25}`)
		expected := map[string]interface{}{
			"name": "Alice",
			"age":  float64(25),
		}
		result, err := JSONToMap(data)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Given invalid JSON, it should return an error", func(t *testing.T) {
		data := []byte(`{name: Alice, age: 25}`)
		_, err := JSONToMap(data)
		if err == nil {
			t.Fatal("Expected an error, but got none")
		}
	})
}

func TestXmlToMap(t *testing.T) {
	t.Run("Given valid XML, it should parse into a map", func(t *testing.T) {
		data := []byte(`<User><Name>Alice</Name><Age>25</Age></User>`)
		expected := map[string]interface{}{
			"User": map[string]interface{}{
				"Name": map[string]interface{}{"#text": "Alice"},
				"Age":  map[string]interface{}{"#text": "25"},
			},
		}
		result, err := XmlToMap(bytes.NewReader(data))
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Given invalid XML, it should return an error", func(t *testing.T) {
		data := []byte(`<User><Name>Alice</Name><Age>25</Age>`) // Invalid XML
		_, err := XmlToMap(bytes.NewReader(data))
		if err == nil {
			t.Fatal("Expected an error, but got none")
		}
	})
}

func TestParseMessage(t *testing.T) {
	t.Run("Given a JSON message, it should parse correctly", func(t *testing.T) {
		data := []byte(`{"name": "Alice", "age": 25}`)
		expected := map[string]interface{}{
			"name": "Alice",
			"age":  float64(25),
		}
		result, err := ParseMessage(data)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Given an XML message, it should parse correctly", func(t *testing.T) {
		data := []byte(`<User><Name>Alice</Name><Age>25</Age></User>`)
		expected := map[string]interface{}{
			"User": map[string]interface{}{
				"Name": map[string]interface{}{"#text": "Alice"},
				"Age":  map[string]interface{}{"#text": "25"},
			},
		}
		result, err := ParseMessage(data)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Given an unknown message format, it should return an error", func(t *testing.T) {
		data := []byte(`name: Alice, age: 25`)
		_, err := ParseMessage(data)
		if err == nil {
			t.Fatal("Expected an error, but got none")
		}
	})
}
