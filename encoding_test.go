package multibayes

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestClassifierJSON(t *testing.T) {
	classifier := NewClassifier()
	classifier.trainWithTestData()

	b, err := classifier.MarshalJSON()
	assert.Equalf(t, nil, err, "Error marshaling JSON: %v\n", err)

	newclass, err := NewClassifierFromJSON(b)
	assert.Equalf(t, nil, err, "Error unmarshaling JSON: %v\n", err)

	assert.Equalf(t, 5, len(newclass.Matrix.Tokens), "Incorrect token length")
	assert.Equalf(t, 2, len(newclass.Matrix.Classes), "Incorrect class length")
}
