package tokens

import (
	//"fmt"
	"regexp"
	"strings"

	//"github.com/blevesearch/bleve/analysis/token_filters/stop_tokens_filter"
	"github.com/blevesearch/bleve/analysis"
	"github.com/blevesearch/bleve/analysis/tokenizers/regexp_tokenizer"
)

type gram struct {
	Tokens [][]byte
}

type TokenizerConf struct {
	regexp    *regexp.Regexp
	NGramSize int64
}

type Tokenizer struct {
	regexp_tokenizer.RegexpTokenizer
	Conf *TokenizerConf
}

func validateConf(tc *TokenizerConf) {
	tc.regexp = regexp.MustCompile(`\w+|\%|\$|\!`)

	if tc.NGramSize == 0 {
		tc.NGramSize = 2
	}
}

func NewTokenizer(tc *TokenizerConf) (*Tokenizer, error) {
	validateConf(tc)

	return &Tokenizer{*regexp_tokenizer.NewRegexpTokenizer(tc.regexp), tc}, nil
}

// Tokenize and Gramify
func (t *Tokenizer) Parse(doc string) []gram {
	// maybe use token types for datetimes or something instead of
	// the actual byte slice
	tokenized := tokensToBytes(t.Tokenize([]byte(strings.ToLower(doc))))
	nTokens := int64(len(tokenized))

	nNGrams := int64(0)
	for i := int64(1); i <= t.Conf.NGramSize; i++ {
		chosen := choose(nTokens, i)
		nNGrams += chosen
	}

	// wowzers
	ngrams := make([]gram, 0, nNGrams)
	for ngramSize := int64(1); ngramSize <= t.Conf.NGramSize; ngramSize++ {
		nNGramsOfSize := choose(nTokens, ngramSize)

		for i := int64(0); i < nNGramsOfSize; i++ {
			ngrams = append(ngrams, gram{tokenized[i:(i + ngramSize)]})
		}
	}

	return ngrams
}

func tokensToBytes(ts analysis.TokenStream) [][]byte {
	bytes := make([][]byte, len(ts))
	for i, token := range ts {
		bytes[i] = token.Term
	}
	return bytes
}

// not a binomial coefficient -- combinations must be sequential
func choose(n, k int64) int64 {
	return max(n-k+int64(1), 0)
}

func max(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}
