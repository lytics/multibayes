package bag

import (
	//"fmt"
	"github.com/drewlanenga/multibayes/tokens"
	"github.com/ryanbressler/CloudForest"
)

type SparseMatrix struct {
	TokenMap map[string]int // map[tokenstring]tokenindex
	Tokens   []map[int]int  // []map[tokenindex]occurence
	ClassMap map[string]int // map[classname]classindex
	Classes  []map[int]int  // []map[classindex]indicator
}

func NewSparseMatrix() *SparseMatrix {
	return &SparseMatrix{
		TokenMap: make(map[string]int),
		Tokens:   make([]map[int]int, 0, 1000), // maybe make bigger
		ClassMap: make(map[string]int),
		Classes:  make([]map[int]int, 0, 1000),
	}
}

func (s *SparseMatrix) Add(ngrams []tokens.NGram, classes []string) {
	s.updateClasses(classes)

	rowTokens := make(map[int]int)
	for _, ngram := range ngrams {
		ngramString := ngram.String()

		tokenIndex, ok := s.TokenMap[ngramString]
		if !ok {
			s.TokenMap[ngramString] = len(s.TokenMap) // + 1
			tokenIndex = s.TokenMap[ngramString]
		}
		rowTokens[tokenIndex]++
	}
	s.Tokens = append(s.Tokens, rowTokens)

	rowClasses := make(map[int]int)
	for _, class := range classes {
		classIndex := s.ClassMap[class]
		rowClasses[classIndex] = 1
	}
	s.Classes = append(s.Classes, rowClasses)
}

// add the class to the class map if it doesn't already exist
func (s *SparseMatrix) updateClasses(classes []string) {
	for _, class := range classes {
		_, ok := s.ClassMap[class]
		if !ok {
			s.ClassMap[class] = len(s.ClassMap)
		}
	}
}

func (s *SparseMatrix) ToFeatureMatrix() map[string]*CloudForest.FeatureMatrix {
	tokenHeaderMap := make(map[string]int, len(s.TokenMap))
	// header names need to have "N:"
	for token, index := range s.TokenMap {
		headerKey := "N:" + token
		tokenHeaderMap[headerKey] = index
	}

	tokenFeatures := ToFeatures(s.TokenMap, s.Tokens, "N:")
	tokenFeatureMatrix := &CloudForest.FeatureMatrix{
		Data:       tokenFeatures,
		Map:        tokenHeaderMap,
		CaseLabels: make([]string, 0),
	}

	classHeaderMap := make(map[string]int, len(s.ClassMap))
	for class, index := range s.ClassMap {
		headerKey := "C:" + class
		classHeaderMap[headerKey] = index
	}
	classFeatures := ToFeatures(s.ClassMap, s.Classes, "C:")
	classFeatureMatrix := &CloudForest.FeatureMatrix{
		Data:       classFeatures,
		Map:        classHeaderMap,
		CaseLabels: make([]string, 0),
	}

	featureMatrices := make(map[string]*CloudForest.FeatureMatrix)
	featureMatrices["tokens"] = tokenFeatureMatrix
	featureMatrices["classes"] = classFeatureMatrix
	return featureMatrices
}

func ToFeatures(tokenMap map[string]int, tokens []map[int]int, prefix string) []CloudForest.Feature {
	tokenFeatures := make([]CloudForest.Feature, len(tokenMap))
	tokenCount := make([]float64, len(tokens))
	// looks like:
	// [ token ]
	// [   0   ]
	// [   1   ] ...
	// iterate over each token
	for tokenname, tokenindex := range tokenMap {
		// iterate over each map contain the locations of tokens
		for rowindex, columnmap := range tokens {
			// if that token doesn't exist, then add a zero to the column
			if _, ok := columnmap[tokenindex]; !ok {
				tokenCount[rowindex] = 0
			} else {
				// if it does, add the count
				tokenCount[rowindex] = float64(columnmap[tokenindex])
			}
		}
		// by the end of this loop, we've created an array of floats with the length
		// equal to the number of rows (number of maps in tokens). basically a column
		// in our feature matrix

		// make a feature here
		f := &CloudForest.DenseNumFeature{
			NumData:    tokenCount,
			Missing:    make([]bool, len(tokens)),
			Name:       prefix + tokenname,
			HasMissing: false,
		}

		// append it to tokenFeatures
		tokenFeatures[tokenindex] = f
	}
	// by the end of this loop we've run through each token
	return tokenFeatures
}
