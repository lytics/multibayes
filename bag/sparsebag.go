package bag

import (
	"fmt"
	"github.com/drewlanenga/multibayes/tokens"
	"github.com/ryanbressler/CloudForest"
)

type SparseMatrix struct {
	Tokens  map[string]*SparseColumn // []map[tokenindex]occurence
	Classes map[string]*SparseColumn // map[classname]classindex
	N       int                      // number of rows currently in the matrix
}

type SparseColumn struct {
	Name string
	Data []int
	N    int
}

func NewSparseColumn(name string) *SparseColumn {
	return &SparseColumn{
		Name: name,
		Data: make([]int, 0, 1000),
		N:    0,
	}
}

func (s *SparseColumn) Add(index int) {
	s.Data = append(s.Data, index)
	s.N++
}

// return the number of rows that contain the column
func (s *SparseColumn) Count() int {
	return len(s.Data)
}

// sparse to dense
func (s *SparseColumn) Expand(n int) []float64 {
	expanded := make([]float64, n)
	for _, index := range s.Data {
		expanded[index] = 1.0
	}
	return expanded
}

func (s *SparseColumn) ToFeature(n int) CloudForest.Feature {
	return &CloudForest.DenseNumFeature{
		NumData:    s.Expand(n),
		Missing:    make([]bool, n), // do we need this if HasMissing is false?
		Name:       s.Name,
		HasMissing: false,
	}
}

func NewSparseMatrix() *SparseMatrix {
	return &SparseMatrix{
		Tokens:  make(map[string]*SparseColumn),
		Classes: make(map[string]*SparseColumn),
		N:       0,
	}
}

func (s *SparseMatrix) Add(ngrams []tokens.NGram, classes []string) {
	for _, class := range classes {
		if _, ok := s.Classes[class]; !ok {
			s.Classes[class] = NewSparseColumn(class)
		}

		s.Classes[class].Add(s.N)
	}

	for _, ngram := range ngrams {
		gramString := ngram.String()
		if _, ok := s.Tokens[gramString]; !ok {
			s.Tokens[gramString] = NewSparseColumn(gramString)
		}

		s.Tokens[gramString].Add(s.N)
	}

	// increment the row counter
	s.N++
}

func (s *SparseMatrix) ToFeatureMatrix() map[string]*CloudForest.FeatureMatrix {
	featureMatrices := make(map[string]*CloudForest.FeatureMatrix)

	featureMatrices["tokens"] = toFeatureMatrix(s.Tokens, s.N)
	featureMatrices["classes"] = toFeatureMatrix(s.Classes, s.N)

	return featureMatrices
}

func toFeatureMatrix(matrixMap map[string]*SparseColumn, nrow int) *CloudForest.FeatureMatrix {
	features := make([]CloudForest.Feature, 0, len(matrixMap))
	featureMap := make(map[string]int)

	i := 0
	for token, column := range matrixMap {
		features = append(features, column.ToFeature(nrow))
		featureMap[token] = i
		i++
	}

	return &CloudForest.FeatureMatrix{
		Data:       features,
		Map:        featureMap,
		CaseLabels: make([]string, 0), // not really relevant here
	}
}
