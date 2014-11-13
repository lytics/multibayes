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
	tc.regexp = regexp.MustCompile(`[0-9A-z_'\-]+|\%|\$`)

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
	alltokens := tokensToStrings(t.Tokenize([]byte(strings.ToLower(doc))))
	filtered := make(map[int][]byte)
	for i, token := range alltokens {
		exclude := false
		for _, stop := range stop_words {
			if token == stop {
				exclude = true
				break
			}
		}

		// possibly check for certain types here (dates|numbers|etc)
		// just stem in the meantime
		token = porterstemmer.StemString(token)

		if !exclude {
			filtered[i] = []byte(token)
		}
	}
	fmt.Printf("\n%+v\n\n^", alltokens)

	for i := 0; i < len(filtered); i++ {
		if token, ok := filtered[i]; ok {
			fmt.Printf("%s ", string(token))
		} else {
			fmt.Printf("* ")
		}
	}
	fmt.Printf("^\n")

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

	for _, ngram := range allNGrams {
		fmt.Println(ngram.String())
	}
	return allNGrams
}

func (t *Tokenizer) tokensToNGrams(tokens [][]byte) []NGram {
	fmt.Printf("\n\ttokensToNGrams:\n\t\tRaw:")
	for _, token := range tokens {
		fmt.Printf(" %s", string(token))
	}
	fmt.Printf("\n")

	nTokens := int64(len(tokens))

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
			ngrams = append(ngrams, NGram{tokens[i:(i + ngramSize)]})
		}
	}

	fmt.Printf("\n\t\tNGrams:")
	for _, ngram := range ngrams {
		fmt.Printf(" %s", ngram.String())
	}
	fmt.Printf("\n")

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
