package multibayes

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestPosterior(t *testing.T) {
	classifier := NewClassifier()
	classifier.MinClassSize = 0
	classifier.trainWithTestData()

	probs := classifier.Posterior("Aaron's dog has tons of fleas")

	assert.Equalf(t, len(classifier.Matrix.Classes), len(probs), "Posterior returned incorrect number of classes")
}
