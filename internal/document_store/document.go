package documentstore

import (
	"encoding/json"
	"fmt"
)

type DocumentFieldType string

const (
	DocumentFieldTypeString DocumentFieldType = "string"
	DocumentFieldTypeNumber DocumentFieldType = "number"
	DocumentFieldTypeBool   DocumentFieldType = "bool"
	DocumentFieldTypeArray  DocumentFieldType = "array"
	DocumentFieldTypeObject DocumentFieldType = "object"
)

type DocumentField struct {
	Type  DocumentFieldType
	Value interface{}
}

type Document struct {
	Fields map[string]DocumentField
}

func MarshalDocument(input any) (*Document, error) {
	data, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input to JSON: %w", err)
	}

	var tempMap map[string]interface{}
	if err := json.Unmarshal(data, &tempMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON to map: %w", err)
	}

	fields := make(map[string]DocumentField)
	for key, value := range tempMap {
		var fieldType DocumentFieldType

		switch value.(type) {
		case string:
			fieldType = DocumentFieldTypeString
		case float64:
			fieldType = DocumentFieldTypeNumber
		case bool:
			fieldType = DocumentFieldTypeBool
		case []interface{}:
			fieldType = DocumentFieldTypeArray
		case map[string]interface{}:
			fieldType = DocumentFieldTypeObject
		default:
			fieldType = DocumentFieldTypeString
		}

		fields[key] = DocumentField{
			Type:  fieldType,
			Value: value,
		}
	}

	doc := Document{
		Fields: fields,
	}

	return &doc, nil
}

func UnmarshalDocument(doc *Document, output any) error {
	tempMap := make(map[string]interface{})
	for key, field := range doc.Fields {
		tempMap[key] = field.Value
	}

	data, err := json.Marshal(tempMap)
	if err != nil {
		return fmt.Errorf("failed to marshal document to JSON: %w", err)
	}

	if err := json.Unmarshal(data, output); err != nil {
		return fmt.Errorf("failed to unmarshal JSON to output: %w", err)
	}
	return nil
}
