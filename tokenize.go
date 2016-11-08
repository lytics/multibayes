package multibayes

import (
	"bytes"
	"encoding/base64"
	"regexp"
	"strings"

	"github.com/blevesearch/bleve/analysis"
	regexp_tokenizer "github.com/blevesearch/bleve/analysis/tokenizer/regexp"
	"github.com/blevesearch/go-porterstemmer"
)

const (
	tokenSeparator = "_"
)

type ngram struct {
	Tokens [][]byte
}

// encodes in base64 for safe comparison
func (ng *ngram) String() string {
	encoded := make([]string, len(ng.Tokens))

	for i, token := range ng.Tokens {
		encoded[i] = string(token)
		//encoded[i] = base64.StdEncoding.EncodeToString(token) // safer?
	}

	return strings.Join(encoded, tokenSeparator)
}

func decodeNGram(s string) (*ngram, error) {
	encodedTokens := strings.Split(s, tokenSeparator)

	tokens := make([][]byte, len(encodedTokens))

	var err error
	for i, encodedToken := range encodedTokens {
		tokens[i], err = base64.StdEncoding.DecodeString(encodedToken)
		if err != nil {
			return nil, err
		}
	}
	return &ngram{tokens}, nil
}

type tokenizerConf struct {
	regexp    *regexp.Regexp
	NGramSize int64
}

type tokenizer struct {
	regexp_tokenizer.RegexpTokenizer
	Conf *tokenizerConf
}

func validateConf(tc *tokenizerConf) {
	tc.regexp = regexp.MustCompile(`[0-9A-z_'\-]+|\%|\$`)

	// TODO: We force NGramSize = 1 so as to create disjoint ngrams,
	// which is necessary for the naive assumption of conditional
	// independence among tokens. It would be great to allow ngrams
	// to be greater than 1 and select only disjoint ngrams from the
	// tokenizer.
	tc.NGramSize = 1
}

func newTokenizer(tc *tokenizerConf) (*tokenizer, error) {
	validateConf(tc)

	return &tokenizer{*regexp_tokenizer.NewRegexpTokenizer(tc.regexp), tc}, nil
}

// Tokenize and Gramify
func (t *tokenizer) Parse(doc string) []ngram {
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
	allNGrams := make([]ngram, 0, 100)
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

func (t *tokenizer) tokensToNGrams(tokens [][]byte) []ngram {
	nTokens := int64(len(tokens))

	nNGrams := int64(0)
	for i := int64(1); i <= t.Conf.NGramSize; i++ {
		chosen := choose(nTokens, i)
		nNGrams += chosen
	}

	ngrams := make([]ngram, 0, nNGrams)
	for ngramSize := int64(1); ngramSize <= t.Conf.NGramSize; ngramSize++ {
		nNGramsOfSize := choose(nTokens, ngramSize)

		for i := int64(0); i < nNGramsOfSize; i++ {
			ngrams = append(ngrams, ngram{tokens[i:(i + ngramSize)]})
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
