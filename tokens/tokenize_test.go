package tokens

import (
	//"fmt"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/drewlanenga/multibayes/testutil"
)

func TestTokenizer(t *testing.T) {
	/*
		docs := []string{
			`This is a sentence with 30% off for $50!`,
			`This is another sentence from 1941-12-07.`,
			`Let's get 50% off together!`,
		}
	*/
	testdata := testutil.GetTestData()

	tokenizer, err := NewTokenizer(&TokenizerConf{
		NGramSize: 2,
	})

	assert.Equalf(t, nil, err, "Error creating tokenizer: %v", err)

	for _, doc := range testdata.Docs {
		//for _, doc := range docs {
		_ = tokenizer.Parse(doc)
		/*
			for _, gram := range grams {
				fmt.Println(gram.String())
			}
		*/
	}

	// test token length here later
}
