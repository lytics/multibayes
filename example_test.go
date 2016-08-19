package multibayes

import (
	"fmt"
)

func Example() {
	documents := []struct {
		Text    string
		Classes []string
	}{
		{
			Text:    "My dog has fleas.",
			Classes: []string{"vet"},
		},
		{
			Text:    "My cat has ebola.",
			Classes: []string{"vet", "cdc"},
		},
		{
			Text:    "Aaron has ebola.",
			Classes: []string{"cdc"},
		},
	}

	classifier := NewClassifier()
	classifier.MinClassSize = 0

	// train the classifier
	for _, document := range documents {
		classifier.Add(document.Text, document.Classes)
	}

	// predict new classes
	probs := classifier.Posterior("Aaron's dog has fleas.")
	fmt.Printf("Posterior Probabilities: vet: %.4f, cdc: %.4f\n", probs["vet"], probs["cdc"])

	// Output: Posterior Probabilities: vet: 0.8571, cdc: 0.2727
}
