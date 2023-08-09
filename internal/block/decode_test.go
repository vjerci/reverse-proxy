package block_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/vjerci/reverse-proxy/internal/block"
)

func TestInterfaceGuardDecoder(t *testing.T) {
	decoder := &block.InterfaceGuardDecoder{}

	var testCases = []struct {
		testName       string
		input          interface{}
		expectedResult interface{}
	}{
		{
			testName: "header_guard",
			input: map[string]string{
				"header": "header",
				"value":  "value",
			},
			expectedResult: &block.HeaderGuard{},
		},
		{
			testName: "query_param_guard",
			input: map[string]string{
				"query_param": "test",
				"value":       "test",
			},
			expectedResult: &block.QueryParamGuard{},
		},
		{
			testName: "method_guard",
			input: map[string]string{
				"method": "get",
			},
			expectedResult: &block.MethodGuard{},
		},
		{
			testName: "path_guard",
			input: map[string]string{
				"path": "/api",
			},
			expectedResult: &block.PathGuard{},
		},
	}

	for _, test := range testCases {
		resultGuard, err := decoder.Decode(test.input)
		if err != nil {
			t.Fatalf("for test %s got %s", test.testName, err)
		}

		if !resultGuard.IsValid() {
			t.Fatalf("for test %s got invalid guard %#v", test.testName, resultGuard)
		}

		if reflect.TypeOf(resultGuard) != reflect.TypeOf(test.expectedResult) {
			t.Fatalf("for test %s, expected type %T but got %T", test.testName, test.expectedResult, resultGuard)
		}
	}
}

func TestInterfaceGuardDecoderError(t *testing.T) {
	decoder := &block.InterfaceGuardDecoder{}

	value, err := decoder.Decode(map[string]string{
		"test": "test",
	})

	if value != nil {
		t.Fatalf("expected nil value but got %#v instead", value)
	}

	if !errors.Is(err, block.ErrDecodeGuard) {
		t.Fatalf("expected ErrFailedToDecodeGuard")
	}
}
