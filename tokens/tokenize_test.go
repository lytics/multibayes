package tokens

import (
	"fmt"
	"testing"

	"github.com/bmizerany/assert"
)

func TestTokenizer(t *testing.T) {
	docs := []string{
		`This is a sentence with 30% off for $50!`,
		`This is another sentence from 1941-12-07.`,
	}

	tokenizer, err := NewTokenizer(&TokenizerConf{
		NGramSize: 2,
	})

	assert.Equalf(t, nil, err, "Error creating tokenizer: %v", err)

	for _, doc := range docs {
		fmt.Println(tokenizer.Parse(doc))
	}
}
