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

	for _, document := range documents {
		classifier.Add(document.Text, document.Classes)
	}

	// predict new classes
	probs := classifier.Posterior("Aaron's dog has fleas.")
	fmt.Printf("Posterior Probabilities: %+v\n", probs)

	// Posterior Probabilities: map[vet:0.8571 cdc:0.2727]
}
