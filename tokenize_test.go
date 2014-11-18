package multibayes

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestTokenizer(t *testing.T) {
	testdata := getTestData()

	tokenize, err := newTokenizer(&tokenizerConf{
		NGramSize: 1,
	})

	assert.Equalf(t, nil, err, "Error creating tokenizer: %v", err)

	for _, doc := range testdata {
		_ = tokenize.Parse(doc.Text)
	}
}
