package bag

import (
	"github.com/drewlanenga/multibayes/tokens"
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
