package multibayes

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestSparseMatrix(t *testing.T) {
	testdata := getTestData()
	tokenizer, err := newTokenizer(&tokenizerConf{
		NGramSize: 1,
	})
	assert.Equalf(t, err, nil, "Error creating new tokenizer")

	sparse := newSparseMatrix()

	for _, document := range testdata {
		ngrams := tokenizer.Parse(document.Text)
		sparse.Add(ngrams, document.Classes)
	}
}
