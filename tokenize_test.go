package multibayes

import (
	"testing"

	"github.com/bmizerany/assert"
	"github.com/drewlanenga/multibayes/testutil"
)

func TestTokenizer(t *testing.T) {
	testdata := testutil.GetTestData()

	tokenizer, err := NewTokenizer(&TokenizerConf{
		NGramSize: 1,
	})

	assert.Equalf(t, nil, err, "Error creating tokenizer: %v", err)

	for _, doc := range testdata.Docs {
		_ = tokenizer.Parse(doc)
	}

	// test token length here later?
}
