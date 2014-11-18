package multibayes

type SparseMatrix struct {
	Tokens  map[string]*SparseColumn `json:"tokens"`  // []map[tokenindex]occurence
	Classes map[string]*SparseColumn `json:"classes"` // map[classname]classindex
	N       int                      `json:"n"`       // number of rows currently in the matrix
}

type SparseColumn struct {
	Data []int `json:"data"`
}

func NewSparseColumn() *SparseColumn {
	return &SparseColumn{
		Data: make([]int, 0, 1000),
	}
}

func (s *SparseColumn) Add(index int) {
	s.Data = append(s.Data, index)
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

func NewSparseMatrix() *SparseMatrix {
	return &SparseMatrix{
		Tokens:  make(map[string]*SparseColumn),
		Classes: make(map[string]*SparseColumn),
		N:       0,
	}
}

func (s *SparseMatrix) Add(ngrams []NGram, classes []string) {
	if len(ngrams) == 0 || len(classes) == 0 {
		return
	}
	for _, class := range classes {
		if _, ok := s.Classes[class]; !ok {
			s.Classes[class] = NewSparseColumn()
		}

		s.Classes[class].Add(s.N)
	}

	for _, ngram := range ngrams {
		gramString := ngram.String()
		if _, ok := s.Tokens[gramString]; !ok {
			s.Tokens[gramString] = NewSparseColumn()
		}

		s.Tokens[gramString].Add(s.N)
	}
	// increment the row counter
	s.N++
}
