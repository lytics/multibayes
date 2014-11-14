package tokens

import (
	"fmt"
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
		fmt.Printf("\nDoc:%s\n", doc)
		grams := tokenizer.Parse(doc)

		for _, gram := range grams {
			fmt.Printf("\tGRAM: %s\n", gram.String())
		}

	}

	// test token length here later
}
