package bag

import (
	//"fmt"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/drewlanenga/multibayes/testutil"
	"github.com/drewlanenga/multibayes/tokens"
)

func TestSparseBag(t *testing.T) {
	testdata := testutil.GetTestData()
	tokenizer, err := tokens.NewTokenizer(&tokens.TokenizerConf{
		NGramSize: 1,
	})
	assert.Equalf(t, err, nil, "Error creating new tokenizer")

	sparse := NewSparseMatrix()

	for i, _ := range testdata.Docs {
		ngrams := tokenizer.Parse(testdata.Docs[i])

		sparse.Add(ngrams, testdata.Classes[i])
	}
}

func TestToFeatureMatrix(t *testing.T) {
	testdata := testutil.GetTestData()
	tokenizer, err := tokens.NewTokenizer(&tokens.TokenizerConf{
		NGramSize: 1,
	})
	assert.Equalf(t, err, nil, "Error creating new tokenizer")

	sparse := NewSparseMatrix()

	for i, _ := range testdata.Docs {
		ngrams := tokenizer.Parse(testdata.Docs[i])

		sparse.Add(ngrams, testdata.Classes[i])
	}

	matrices := sparse.ToFeatureMatrix()

	assert.Equal(t, len(sparse.Classes), len(matrices["classes"].Data), "Wrong length")
	assert.Equal(t, len(sparse.Tokens), len(matrices["tokens"].Data), "Wrong length")
}
