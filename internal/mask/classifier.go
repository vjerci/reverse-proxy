package mask

import (
	"regexp"
)

type Classifier interface {
	ClassifyField(fieldName string, fieldType FieldType) (isClassified bool)
}

type PIIClassifier struct {
	patterns []PIIPattern
}

func NewPIIClassifier(patterns []PIIPattern) *PIIClassifier {
	return &PIIClassifier{
		patterns: patterns,
	}
}

func (classifier *PIIClassifier) ClassifyField(fieldName string, fieldType FieldType) bool {
	for _, pattern := range classifier.patterns {
		if pattern.IsPII([]byte(fieldName)) {
			return true
		}
	}

	return false
}

type PIIPattern interface {
	IsPII(input []byte) bool
}

type PIIClassifierPattern struct {
	Regexp *regexp.Regexp
}

func (pattern *PIIClassifierPattern) IsPII(input []byte) bool {
	return pattern.Regexp.Match(input)
}

func NewDefaultPIIPatterns() []PIIPattern {
	// I wanted to use https://github.com/Bearer/bearer/tree/main/pkg/classification but it is not exactly extensible and simple enough to use in demo project
	return []PIIPattern{
		&PIIClassifierPattern{
			Regexp: regexp.MustCompile(`\w*email\w*`),
		},
		&PIIClassifierPattern{
			Regexp: regexp.MustCompile(`\w*name\w*`),
		},
		&PIIClassifierPattern{
			Regexp: regexp.MustCompile(`\w*gender\w*`),
		},
	}
}
