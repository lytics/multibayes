package multibayes

import (
	"testing"

	//"github.com/bmizerany/assert"
)

func TestPosterior(t *testing.T) {
	classifier := NewClassifier()

	testdata := getTestData()
	for _, document := range testdata {
		classifier.Add(document.Text, document.Classes)
	}

	for _, document := range testdata {
		_ = classifier.Posterior(document.Text)
	}
}
