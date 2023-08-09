package mask

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

type FieldType string

var FieldTypeString = FieldType("string")
var FieldTypeBool = FieldType("bool")
var FieldTypeFloat64 = FieldType("float64")

var ErrDecodeJSON = errors.New("failed to decode json")

type Inspector interface {
	Inspect(bytes []byte) ([]byte, error)
}

type JSONInspector struct {
	masker     Mask
	classifier Classifier
}

func NewJSONInspector(mask Mask, classifier Classifier) Inspector {
	return &JSONInspector{
		masker:     mask,
		classifier: classifier,
	}
}

func (inspector *JSONInspector) Inspect(input []byte) ([]byte, error) {
	var data interface{}
	err := json.NewDecoder(bytes.NewBuffer(input)).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDecodeJSON, err)
	}

	data = inspector.inspectInput(data)

	buff := bytes.NewBuffer(nil)
	json.NewEncoder(buff).Encode(&data)

	return buff.Bytes(), nil
}

func (inspector *JSONInspector) inspectInput(input interface{}) interface{} {
	switch inputTyped := input.(type) {
	case map[string]interface{}:
		return inspector.inspectMap(inputTyped)
	case []interface{}:
		return inspector.inspectArray(inputTyped)
	default:
		return input
	}
}

func (inspector *JSONInspector) inspectMap(input map[string]interface{}) map[string]interface{} {
	for key, value := range input {
		switch value.(type) {
		case float64:
			if inspector.classifier.ClassifyField(key, FieldTypeFloat64) {
				input[key] = inspector.masker.Mask(value, FieldTypeFloat64)
			}
		case string:
			if inspector.classifier.ClassifyField(key, FieldTypeString) {
				input[key] = inspector.masker.Mask(value, FieldTypeString)
			}
		case bool:
			if inspector.classifier.ClassifyField(key, FieldTypeBool) {
				input[key] = inspector.masker.Mask(value, FieldTypeBool)
			}
		default:
			input[key] = inspector.inspectInput(value)
		}
	}

	return input
}

func (inspector *JSONInspector) inspectArray(input []interface{}) []interface{} {
	for key, value := range input {
		input[key] = inspector.inspectInput(value)
	}

	return input
}
