package tokens

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"

	"github.com/blevesearch/bleve/analysis"
	//"github.com/blevesearch/bleve/analysis/token_filters/stop_tokens_filter"
	"github.com/blevesearch/bleve/analysis/tokenizers/regexp_tokenizer"
	"github.com/blevesearch/go-porterstemmer"
	//"github.com/lytics/lio/src/catalog/prob"
	//"github.com/lytics/lio/src/catalog/prob/data"
)

const (
	tokenSeparator = "_"
)

var (
	stop_words = []string{`i`, `me`, `my`, `myself`, `we`, `our`, `ours`, `ourselves`, `you`, `your`, `yours`, `yourself`, `yourselves`, `he`, `him`, `his`, `himself`, `she`, `her`, `hers`, `herself`, `it`, `its`, `itself`, `they`, `them`, `their`, `theirs`, `themselves`, `what`, `which`, `who`, `whom`, `this`, `that`, `these`, `those`, `am`, `is`, `are`, `was`, `were`, `be`, `been`, `being`, `have`, `has`, `had`, `having`, `do`, `does`, `did`, `doing`, `would`, `should`, `could`, `ought`, `i'm`, `you're`, `he's`, `she's`, `it's`, `we're`, `they're`, `i've`, `you've`, `we've`, `they've`, `i'd`, `you'd`, `he'd`, `she'd`, `we'd`, `they'd`, `i'll`, `you'll`, `he'll`, `she'll`, `we'll`, `they'll`, `isn't`, `aren't`, `wasn't`, `weren't`, `hasn't`, `haven't`, `hadn't`, `doesn't`, `don't`, `didn't`, `won't`, `wouldn't`, `shan't`, `shouldn't`, `can't`, `cannot`, `couldn't`, `mustn't`, `let's`, `that's`, `who's`, `what's`, `here's`, `there's`, `when's`, `where's`, `why's`, `how's`, `a`, `an`, `the`, `and`, `but`, `if`, `or`, `because`, `as`, `until`, `while`, `of`, `at`, `by`, `for`, `with`, `about`, `against`, `between`, `into`, `through`, `during`, `before`, `after`, `above`, `below`, `to`, `from`, `up`, `down`, `in`, `out`, `on`, `off`, `over`, `under`, `again`, `further`, `then`, `once`, `here`, `there`, `when`, `where`, `why`, `how`, `all`, `any`, `both`, `each`, `few`, `more`, `most`, `other`, `some`, `such`, `no`, `nor`, `not`, `only`, `own`, `same`, `so`, `than`, `too`, `very`}
)

type NGram struct {
	Tokens [][]byte
}

// encodes in base64 for safe comparison
func (ng *NGram) String() string {
	encoded := make([]string, len(ng.Tokens))

	for i, token := range ng.Tokens {
		encoded[i] = string(token)
		//encoded[i] = base64.StdEncoding.EncodeToString(token)
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
	tc.regexp = regexp.MustCompile(`\w+|\%|\$`)

	if tc.NGramSize == 0 {
		tc.NGramSize = 2
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
	tokenizedstr := tokensToStrings(t.Tokenize([]byte(strings.ToLower(doc))))
	var filtered []string
	var exists bool
	for _, token := range tokenizedstr {
		for _, stop := range stop_words {
			if token == stop {
				exists = true
				break
			}
		}
		if !exists {
			filtered = append(filtered, token)
		}
		exists = false
	}
	//fmt.Println(tokenizedstr)
	//fmt.Println(filtered)

	tokenized := make([][]byte, len(filtered))
	for i, token := range filtered {
		stem := porterstemmer.StemString(token)
		tokenizedstr[i] = stem
		tokenized[i] = []byte(stem)
	}
	fmt.Println(tokenizedstr)

	nTokens := int64(len(tokenized))

	nNGrams := int64(0)
	for i := int64(1); i <= t.Conf.NGramSize; i++ {
		chosen := choose(nTokens, i)
		nNGrams += chosen
	}

	// wowzers
	ngrams := make([]NGram, 0, nNGrams)
	for ngramSize := int64(1); ngramSize <= t.Conf.NGramSize; ngramSize++ {
		nNGramsOfSize := choose(nTokens, ngramSize)

		for i := int64(0); i < nNGramsOfSize; i++ {
			ngrams = append(ngrams, NGram{tokenized[i:(i + ngramSize)]})
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

func tokensToStrings(ts analysis.TokenStream) []string {
	strs := make([]string, len(ts))
	for i, token := range ts {
		strs[i] = string(token.Term)
	}
	return strs
}

func tokensToStemmedBytes(ts analysis.TokenStream) [][]byte {
	bytes := make([][]byte, len(ts))
	for i, token := range ts {
		stem := porterstemmer.StemString(string(token.Term))
		bytes[i] = []byte(stem)
	}
	return bytes
}

func tokensToStemmedStrings(ts analysis.TokenStream) []string {
	strs := make([]string, len(ts))
	for i, token := range ts {
		strs[i] = porterstemmer.StemString(string(token.Term))
	}
	return strs
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
