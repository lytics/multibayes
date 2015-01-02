package multibayes

type sparseMatrix struct {
	Tokens  map[string]*sparseColumn `json:"tokens"`  // []map[tokenindex]occurence
	Classes map[string]*sparseColumn `json:"classes"` // map[classname]classindex
	N       int                      `json:"n"`       // number of rows currently in the matrix
}

type sparseColumn struct {
	Data []int `json:"data"`
}

func newSparseColumn() *sparseColumn {
	return &sparseColumn{
		Data: make([]int, 0, 1000),
	}
}

func (s *sparseColumn) Add(index int) {
	s.Data = append(s.Data, index)
}

// return the number of rows that contain the column
func (s *sparseColumn) Count() int {
	return len(s.Data)
}

// sparse to dense
func (s *sparseColumn) Expand(n int) []float64 {
	expanded := make([]float64, n)
	for _, index := range s.Data {
		expanded[index] = 1.0
	}
	return expanded
}

func newSparseMatrix() *sparseMatrix {
	return &sparseMatrix{
		Tokens:  make(map[string]*sparseColumn),
		Classes: make(map[string]*sparseColumn),
		N:       0,
	}
}

func (s *sparseMatrix) Add(ngrams []ngram, classes []string) {
	if len(ngrams) == 0 || len(classes) == 0 {
		return
	}
	for _, class := range classes {
		if _, ok := s.Classes[class]; !ok {
			s.Classes[class] = newSparseColumn()
		}

		s.Classes[class].Add(s.N)
	}

	// add ngrams uniquely
	added := make(map[string]int)
	for _, ngram := range ngrams {
		gramString := ngram.String()
		if _, ok := s.Tokens[gramString]; !ok {
			s.Tokens[gramString] = newSparseColumn()
		}

		// only add the document index once for the ngram
		if _, ok := added[gramString]; !ok {
			added[gramString] = 1
			s.Tokens[gramString].Add(s.N)
		}
	}
	// increment the row counter
	s.N++
}
