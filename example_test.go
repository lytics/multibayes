package multibayes

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
}
