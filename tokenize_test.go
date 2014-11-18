package multibayes

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestTokenizer(t *testing.T) {
	testdata := getTestData()

	tokenizer, err := NewTokenizer(&TokenizerConf{
		NGramSize: 1,
	})

	assert.Equalf(t, nil, err, "Error creating tokenizer: %v", err)

	for _, doc := range testdata {
		_ = tokenizer.Parse(doc.Text)
	}

	// test token length here later?
}
