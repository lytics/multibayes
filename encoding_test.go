package multibayes

import (
	"fmt"
	"github.com/bmizerany/assert"
	"testing"
)

func TestClassifierJSON(t *testing.T) {
	classifier := NewClassifier()

	testdata := getTestData()

	for _, document := range testdata {
		classifier.Add(document.Text, document.Classes)
	}

	b, err := classifier.MarshalJSON()
	assert.Equalf(t, nil, err, "Error marshaling JSON: %v\n", err)

	fmt.Println(string(b))

	newclass, err := NewClassifierFromJSON(b)
	assert.Equalf(t, nil, err, "Error unmarshaling JSON: %v\n", err)

	fmt.Println(newclass)
}
