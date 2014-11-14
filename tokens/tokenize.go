package tokens

import (
	"bytes"
	"encoding/base64"
	"regexp"
	"strings"

	"github.com/blevesearch/bleve/analysis"
	"github.com/blevesearch/bleve/analysis/tokenizers/regexp_tokenizer"
	"github.com/blevesearch/go-porterstemmer"
)

const (
	tokenSeparator = "_"
)

type NGram struct {
	Tokens [][]byte
}

// encodes in base64 for safe comparison
func (ng *NGram) String() string {
	encoded := make([]string, len(ng.Tokens))

	for i, token := range ng.Tokens {
		encoded[i] = string(token)
		//encoded[i] = base64.StdEncoding.EncodeToString(token) // safer?
	}

	return strings.Join(encoded, tokenSeparator)
}

func DecodeNGram(s string) (*NGram, error) {
	encodedTokens := strings.Split(s, tokenSeparator)

	tokens := make([][]byte, len(encodedTokens))

	var err error
	for i, encodedToken := range encodedTokens {
		tokens[i], err = base64.StdEncoding.DecodeString(encodedToken)
		if err != nil {
			return nil, err
		}
	}
	return &NGram{tokens}, nil
}

type TokenizerConf struct {
	regexp    *regexp.Regexp
	NGramSize int64
}

type Tokenizer struct {
	regexp_tokenizer.RegexpTokenizer
	Conf *TokenizerConf
}

type StopFilter struct {
}

func validateConf(tc *TokenizerConf) {
	tc.regexp = regexp.MustCompile(`[0-9A-z_'\-]+|\%|\$`)

	if tc.NGramSize == 0 {
		tc.NGramSize = 1
	}
}

func NewTokenizer(tc *TokenizerConf) (*Tokenizer, error) {
	validateConf(tc)

	return &Tokenizer{*regexp_tokenizer.NewRegexpTokenizer(tc.regexp), tc}, nil
}

// Tokenize and Gramify
func (t *Tokenizer) Parse(doc string) []NGram {
	// maybe use token types for datetimes or something instead of
	// the actual byte slice
	alltokens := t.Tokenize([]byte(strings.ToLower(doc)))
	filtered := make(map[int][]byte)
	for i, token := range alltokens {
		exclude := false
		for _, stop := range stopbytes {
			if bytes.Equal(token.Term, stop) {
				exclude = true
				break
			}
		}

		if exclude {
			continue
		}

		tokenString := porterstemmer.StemString(string(token.Term))
		//tokenBytes := porterstemmer.Stem(token.Term) // takes runes, not bytes

		if token.Type == analysis.Numeric {
			tokenString = "NUMBER"
		} else if token.Type == analysis.DateTime {
			tokenString = "DATE"
		}

		filtered[i] = []byte(tokenString)
	}

	// only consider sequential terms as candidates for ngrams
	// terms separated by stopwords are ineligible
	allNGrams := make([]NGram, 0, 100)
	currentTokens := make([][]byte, 0, 100)

	lastObserved := -1
	for i, token := range filtered {
		if (i - 1) != lastObserved {

			ngrams := t.tokensToNGrams(currentTokens)
			allNGrams = append(allNGrams, ngrams...)

			currentTokens = make([][]byte, 0, 100)
		}

		currentTokens = append(currentTokens, token)
		lastObserved = i
	}

	// bring in the last one
	if len(currentTokens) > 0 {
		ngrams := t.tokensToNGrams(currentTokens)
		allNGrams = append(allNGrams, ngrams...)
	}

	return allNGrams
}

func (t *Tokenizer) tokensToNGrams(tokens [][]byte) []NGram {
	nTokens := int64(len(tokens))

	nNGrams := int64(0)
	for i := int64(1); i <= t.Conf.NGramSize; i++ {
		chosen := choose(nTokens, i)
		nNGrams += chosen
	}

	ngrams := make([]NGram, 0, nNGrams)
	for ngramSize := int64(1); ngramSize <= t.Conf.NGramSize; ngramSize++ {
		nNGramsOfSize := choose(nTokens, ngramSize)

		for i := int64(0); i < nNGramsOfSize; i++ {
			ngrams = append(ngrams, NGram{tokens[i:(i + ngramSize)]})
		}
	}

	return ngrams
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
