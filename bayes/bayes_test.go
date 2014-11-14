package bayes

import (
	"fmt"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/drewlanenga/multibayes/bag"
	"github.com/drewlanenga/multibayes/testutil"
	"github.com/drewlanenga/multibayes/tokens"
)

func TestPosterior(t *testing.T) {
	testdata := testutil.GetTestData()
	tokenizer, err := tokens.NewTokenizer(&tokens.TokenizerConf{
		NGramSize: 1,
	})
	assert.Equalf(t, err, nil, "Error creating new tokenizer")

	sparse := bag.NewSparseMatrix()
	for i, _ := range testdata.Docs {
		ngrams := tokenizer.Parse(testdata.Docs[i])
		sparse.Add(ngrams, testdata.Classes[i])
	}

	for i, _ := range testdata.Docs {
		fmt.Printf("Subject: %s\n", testdata.Docs[i])
		subject := tokenizer.Parse(testdata.Docs[i])
		predictions := Posterior(sparse, subject)
		for class, pred := range predictions {
			if pred > 0.1 {
				fmt.Printf("%s: %.4f,\n", class, pred)
			}
		}
		fmt.Println()
	}

}
