Multibayes
==========

Multi-multinomial naive Bayesian document classification.

Often in document classification, a document may have more than one relevant classification -- a question on [stackoverflow](http://stackoverflow.com) might have tags "go", "map", and "interface".  The multibayes library strives to offer efficient storage and calculation of Bayesian posterior classification probabilities.

## Example

```go
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
```
