package config_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vjerci/reverse-proxy/internal/config"
)

func TestLoadConfigErrors(t *testing.T) {
	testCases := []struct {
		testName      string
		input         string
		expectedError error
	}{
		{
			testName:      "empty_file_path",
			input:         "",
			expectedError: config.ErrConfigNotSet,
		},
		{
			testName:      "non_existing_file",
			input:         "./testdata/non_existing",
			expectedError: config.ErrConfigRead,
		},
		{
			testName:      "empty_file_path",
			input:         "./testdata/faulty_config.json",
			expectedError: config.ErrConfigJSON,
		},
	}

	for _, test := range testCases {
		config, err := config.Load(test.input)

		testName := t.Name() + "-" + test.testName

		t.Logf("Current test filename: %s", test.input)

		assert.Nil(t, config, "%s for test expected config to be nil", testName)

		if !errors.Is(err, test.expectedError) {
			t.Fatalf("expected to get error %s got %s instead", test.expectedError, err)
		}
	}
}

func TestLoadConfigSuccess(t *testing.T) {
	testCases := []struct {
		forwardHost   string
		forwardScheme string
		input         string
		block         [][]interface{}
	}{
		{
			forwardHost:   "localhost:8000",
			forwardScheme: "http",
			input:         "./testdata/config.json",
			block: [][]interface{}{
				{

					map[string]interface{}{
						"method": "get",
					},
					map[string]interface{}{
						"query_param": "test",
						"value":       "test",
					},
				},
			},
		},
	}

	for _, test := range testCases {
		config, err := config.Load(test.input)

		assert.Nil(t, err, "expected err to be nil")

		assert.Equal(t, test.forwardHost, config.ForwardHost, "expected forward host to be '%s' got '%s' instead", test.forwardHost, config.ForwardHost)

		assert.Equal(t, test.forwardScheme, config.ForwardScheme, "expected forward scheme to be '%s' got '%s' instead", test.forwardScheme, config.ForwardScheme)

		assert.Equal(t, test.block, config.Block, "expected block to be %#v, got %#v instead", test.block, config.Block)
	}
}
