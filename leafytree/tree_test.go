package leafytree

import (
	"fmt"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/drewlanenga/multibayes/testutil"
	"github.com/drewlanenga/multibayes/tokens"
)

func TestTree(t *testing.T) {
	testdata := testutil.GetTestData()

	tokenizer, err := tokens.NewTokenizer(&tokens.TokenizerConf{
		NGramSize: 1,
	})
	assert.Equalf(t, nil, err, "Error creating tokenizer: %v", err)

	tree := NewTree()

	for i, doc := range testdata.Docs {
		ngrams := tokenizer.Parse(doc)
		tree.Learn(ngrams, testdata.Classes[i])
	}

	fmt.Printf("\n\n-------------------------------\n\n")
	for _, doc := range testdata.Docs {
		fmt.Printf("\n\n\nDoc: %s\n", doc)
		predicted := tree.Predict(tokenizer.Parse(doc))

		for key, prob := range predicted {
			if prob > 0.2 {
				fmt.Printf("\n\t%s: %f", key, prob)
			}
		}
	}
}
