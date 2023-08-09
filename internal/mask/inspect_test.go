package mask

import (
	"errors"
	"testing"

	"github.com/bradleyjkemp/cupaloy/v2"
)

type ClassifierMock struct{}

func (classifier *ClassifierMock) ClassifyField(fieldName string, fieldType FieldType) bool {
	return true
}

func TestJSONInspector(t *testing.T) {
	jsonInspector := NewJSONInspector(NewJSONMask(), &ClassifierMock{})

	testCases := []struct {
		testName string
		json     string
	}{
		{
			testName: "nested_object",
			json: `{
				"object": {
					"nestedObject":{
						"string": "string"
					}
				},
				"string": "hello",
				"number": 1,
				"bool": false
			}`,
		},
		{
			testName: "string",
			json:     `"hello"`,
		},
		{
			testName: "boolean",
			json:     `true`,
		},
		{
			testName: "number",
			json:     `0`,
		},
		{
			testName: "array",
			json:     `[0,0,0,0]`,
		},
		{
			testName: "array_of_objects",
			json: `[{
				"string": "string"
			},{
				"bool": false
			},{
				"number": 0.01
			}]`,
		},
	}

	for _, test := range testCases {
		output, err := jsonInspector.Inspect([]byte(test.json))

		if err != nil {
			t.Fatalf("expected to get output from json inspector %s", err)
		}

		snapshotName := (t.Name() + "-" + test.testName)

		err = cupaloy.SnapshotWithName(snapshotName, output)
		if err != nil {
			t.Fatalf("for test %s snapshot doesn't match", snapshotName)
		}
	}
}

func TestJSONInspectorError(t *testing.T) {
	jsonInspector := NewJSONInspector(NewJSONMask(), &ClassifierMock{})

	input := "{[}]}"

	output, err := jsonInspector.Inspect([]byte(input))

	if output != nil {
		t.Fatalf("expected nil output but got one instead")
	}

	if !errors.Is(err, ErrDecodeJSON) {
		t.Fatalf("expected to get wrapped error")
	}
}
