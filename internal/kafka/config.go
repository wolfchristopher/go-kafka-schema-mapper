package kafka

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"reflect"
	"strings"
)

// MapSchema dynamically maps the schema of the input map
func MapSchema(data map[string]interface{}) map[string]interface{} {
	schema := make(map[string]interface{})
	for key, value := range data {
		if value == nil {
			schema[key] = "<nil>"
			continue
		}

		// Store the type of the value
		schema[key] = reflect.TypeOf(value).String()

		// If the value is a nested structure, recursively map its schema
		if reflect.TypeOf(value).Kind() == reflect.Map {
			schema[key] = MapSchema(value.(map[string]interface{}))
		}
	}
	return schema
}

// XmlToMap xmlToMap dynamically converts XML into a map[string]interface{}
func XmlToMap(reader io.Reader) (map[string]interface{}, error) {
	decoder := xml.NewDecoder(reader)
	var stack []map[string]interface{}
	current := make(map[string]interface{})

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break // End of XML
		}
		if err != nil {
			return nil, fmt.Errorf("error decoding XML: %v", err)
		}

		switch tok := token.(type) {
		case xml.StartElement:
			element := make(map[string]interface{})
			for _, attr := range tok.Attr {
				element[attr.Name.Local] = attr.Value
			}

			stack = append(stack, current)

			current = element

		case xml.EndElement:
			if len(stack) == 0 {
				return current, nil
			}
			parent := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if existing, found := parent[tok.Name.Local]; found {
				switch existing := existing.(type) {
				case []interface{}:
					parent[tok.Name.Local] = append(existing, current)
				default:
					parent[tok.Name.Local] = []interface{}{existing, current}
				}
			} else {
				parent[tok.Name.Local] = current
			}

			current = parent

		case xml.CharData:
			content := string(tok)
			if len(content) > 0 {
				if current == nil {
					current = make(map[string]interface{})
				}
				current["#text"] = content
			}
		}
	}

	return current, nil
}

// ParseMessage Detect message format and parse it into a map
func ParseMessage(data []byte) (map[string]interface{}, error) {
	trimmedData := strings.TrimSpace(string(data))

	if strings.HasPrefix(trimmedData, "{") {
		return JSONToMap(data)
	} else if strings.HasPrefix(trimmedData, "<") {
		return XmlToMap(bytes.NewReader(data))
	} else {
		return nil, fmt.Errorf("unknown message format")
	}
}

// JSONToMap Parse JSON into a map[string]interface{}
func JSONToMap(data []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}
	return result, nil
}
