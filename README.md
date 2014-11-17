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

	tokenizer, _ := NewTokenizer(&TokenizerConf{
		NGramSize: 1,
	})

	sparse := NewSparseMatrix()

	for _, document := range documents {
		ngrams := tokenizer.Parse(document.Text)
		sparse.Add(ngrams, document.Classes)
	}

	// predict new classes
	probs := Posterior(sparse, tokenizer.Parse("Aaron's dog has fleas."))
	fmt.Printf("Posterior Probabilities: %+v\n", probs)

	// Posterior Probabilities: map[vet:0.8571 cdc:0.2727]
```

## Notes

The tokenizer tries to create *disjoint* ngrams for an input document.  (If the ngrams aren't disjoint, the naive assumption of conditional independence between tokens is invalid.)  The tokenizer currently doesn't return disjoint ngrams, so to ensure "valid" independence assumptions, make sure that the tokenizer is configured with `NGramSize` of 1.

PRs are happily accepted!
