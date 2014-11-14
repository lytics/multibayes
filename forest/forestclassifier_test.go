package forestclassifier

import (
	"fmt"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/drewlanenga/multibayes/testutil"
	"github.com/drewlanenga/multibayes/tokens"
)

func TestLearn(t *testing.T) {
	testdata := testutil.GetTestData()
	tokenizer, err := tokens.NewTokenizer(&tokens.TokenizerConf{
		NGramSize: 1,
	})
	assert.Equalf(t, err, nil, "Error creating new tokenizer")

	docs := make([][]tokens.NGram, len(testdata.Docs))
	for i, _ := range testdata.Docs {
		docs[i] = tokenizer.Parse(testdata.Docs[i])
	}
}

func TestPredict(t *testing.T) {
	testdata := testutil.GetTestData()
	tokenizer, err := tokens.NewTokenizer(&tokens.TokenizerConf{
		NGramSize: 1,
	})
	assert.Equalf(t, err, nil, "Error creating new tokenizer")

	docs := make([][]tokens.NGram, len(testdata.Docs))
	for i, _ := range testdata.Docs {
		docs[i] = tokenizer.Parse(testdata.Docs[i])
	}
	forestBag, matrices := Learn(docs, testdata.Classes)

	for i, _ := range testdata.Docs {
		//if i > 0 {
		//	break
		//}
		fmt.Printf("Subject line: \t %v\n", testdata.Docs[i])
		predictions := forestBag.Predict(matrices["tokens"], docs[i])
		fmt.Println(predictions)
		fmt.Println()
	}
}
