package mask_test

import (
	"testing"

	"github.com/vjerci/reverse-proxy/internal/mask"
)

func TestDefaultPIIClassification(t *testing.T) {
	patterns := mask.NewDefaultPIIPatterns()
	classifier := mask.NewPIIClassifier(patterns)

	testCases := []struct {
		fieldName      string
		FieldType      mask.FieldType
		expectedResult bool
	}{
		{
			fieldName:      "first_name",
			FieldType:      mask.FieldTypeString,
			expectedResult: true,
		},
		{
			fieldName:      "last_name",
			FieldType:      mask.FieldTypeString,
			expectedResult: true,
		},
		{
			fieldName:      "personal_email",
			FieldType:      mask.FieldTypeString,
			expectedResult: true,
		},
		{
			fieldName:      "gender_identified",
			FieldType:      mask.FieldTypeString,
			expectedResult: true,
		},
		{
			fieldName:      "usage_stats",
			FieldType:      mask.FieldTypeFloat64,
			expectedResult: false,
		},
		{
			fieldName:      "isAdmin",
			FieldType:      mask.FieldTypeBool,
			expectedResult: false,
		},
	}

	for _, testCase := range testCases {
		result := classifier.ClassifyField(testCase.fieldName, testCase.FieldType)

		if result != testCase.expectedResult {
			t.Fatalf("expected to get %t for fieldName %s", testCase.expectedResult, testCase.fieldName)
		}
	}
}
