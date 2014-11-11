package forestclassifier

import (
	"github.com/drewlanenga/matrix"
	"github.com/drewlanenga/tokens"
	"github.com/lio/src/scores/mondrian"
	"github.com/ryanbressler/CloudForest"
)

type ForestBag map[string]*mondrian.MondrianForest
// given subject lines and classes
// tokenize subject lines
// make feature matrix for each type of email
// build forest from feature matrix
// 

func docsToTokens(docs []string) [][]tokens.NGram {
	tokenized := make([][]tokens.NGram, len(docs))
	tokenizer, err := NewTokenizer(&TokenizerConf{
		NGramSize: 1,
	})

	for i, doc := range docs {
		tokenized[i] = tokenizer.Parse(doc)
	}

	return tokenized
}

func tokensToFeatureMatrix(tokens [][]tokens.NGram, classes [][]string) {


	for _, class := range classes {
		featureMatrix := &CloudForest.FeatureMatrix{
			Data: ,
			Map: ,
			CaseLabels: ,
		}
	}
}

/*
func Train(ngrams []tokens.NGram, classes []string) *ForestBag {
	// make Feature Matrix
	features := make([]CloudForest.Feature, len())

	// build ForestBag
}
*/

func Train(doc []string, classes []string) *ForestBag {
	// make Feature Matrix
	features := make([]CloudForest.Feature, len(doc))
	// 		|
	// doc1
	// doc2

}

func (forest *ForestBag) Predict(subject tokens.Ngram) {
	//

}
