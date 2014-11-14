package bayes

import (
	"github.com/drewlanenga/multibayes/matrix"
	"github.com/drewlanenga/multibayes/tokens"
)

func Posterior(s *matrix.SparseMatrix, subject []tokens.NGram) map[string]float64 {
	predictions := make(map[string]float64)

	var prior float64
	likelihood := 1.0
	var conditional float64
	evidence := make(map[string]float64)
	var totalevidence float64

	for class, classcolumn := range s.Classes {
		prior = float64(classcolumn.Count()) / float64(s.N)
		// check if subject token is in our token sparse matrix
		for _, subjecttoken := range subject {
			if tokencolumn, ok := s.Tokens[subjecttoken.String()]; ok {
				// if a subject token is, then check if this token has
				// occurred for this class type
				array := intersection(classcolumn.Data, tokencolumn.Data)
				conditional = float64(len(array))
				// conditional should be the percentage of times this token
				// occurred for a specific class
				conditional = (conditional + 1) / (float64(len(classcolumn.Data)) + 1)
				likelihood *= conditional
				conditional = 0.0
			}
		}
		evidence[class] = prior * likelihood
		totalevidence += prior * likelihood
		likelihood = 1.0
	}
	for class, _ := range s.Classes {
		predictions[class] = evidence[class] / totalevidence
	}
	return predictions
}

func intersection(array1, array2 []int) []int {
	var newarray []int
	for _, elem1 := range array1 {
		for _, elem2 := range array2 {
			if elem1 == elem2 {
				newarray = append(newarray, elem1)
			}
		}
	}
	return newarray
}
