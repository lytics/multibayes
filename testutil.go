package multibayes

type document struct {
	Text    string
	Classes []string
}

func getTestData() []document {

	documents := []document{
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

	return documents
}
