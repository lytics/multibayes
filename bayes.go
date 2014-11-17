package multibayes

import (
	"fmt"
	"github.com/drewlanenga/multibayes/matrix"
	"github.com/drewlanenga/multibayes/tokens"
	"math"
)

var (
	smoother = 1 // laplace
)

func Posterior(s *matrix.SparseMatrix, subject []tokens.NGram) map[string]float64 {
	predictions := make(map[string]float64)

	for class, classcolumn := range s.Classes {
		n := classcolumn.Count()

		priors := []float64{
			float64(n+smoother) / float64(s.N+(smoother*2)),     // P(C=Y)
			float64(s.N-n+smoother) / float64(s.N+(smoother*2)), // P(C=N)
		}

		loglikelihood := []float64{1.0, 1.0}

		// check if subject token is in our token sparse matrix
		for _, subjecttoken := range subject {
			if tokencolumn, ok := s.Tokens[subjecttoken.String()]; ok {
				// conditional probability the token occurs for the class
				joint := intersection(tokencolumn.Data, classcolumn.Data)
				conditional := float64(joint+smoother) / float64(n+(smoother*2)) // P(F|C=Y)
				loglikelihood[0] += math.Log(conditional)

				// conditional probability the token occurs if the class doesn't apply
				not := notintersection(tokencolumn.Data, classcolumn.Data)
				notconditional := float64(not+smoother) / float64(n+(smoother*2)) // P(F|C=N)
				loglikelihood[1] += math.Log(notconditional)
			}
		}

		likelihood := []float64{
			math.Exp(loglikelihood[0]),
			math.Exp(loglikelihood[1]),
		}

		prob := bayesRule(priors, likelihood) // P(C|F)
		predictions[class] = prob[0]
	}

	// just for debugging -- delete later
	fmt.Printf("\n\tPredictions:")
	for class, prob := range predictions {
		if prob > 0.1 {
			fmt.Printf("\n\t\t%s, %.8f", class, prob)
		}
	}
	fmt.Printf("\n\n")

	return predictions
}

func bayesRule(prior, likelihood []float64) []float64 {

	posterior := make([]float64, len(prior))

	sum := 0.0
	for i, _ := range prior {
		combined := prior[i] * likelihood[i]

		posterior[i] = combined
		sum += combined
	}

	// scale the likelihoods
	for i, _ := range posterior {
		posterior[i] /= sum
	}

	return posterior
}

// elements that are in both array1 and array2
func intersection(array1, array2 []int) int {
	var count int
	for _, elem1 := range array1 {
		for _, elem2 := range array2 {
			if elem1 == elem2 {
				count++
				break
			}
		}
	}
	return count
}

// given it's not an element of array2, the count of array1
func notintersection(array1, array2 []int) int {
	var count int

	for _, elem1 := range array1 {
		isElement := false
		for _, elem2 := range array2 {
			if elem1 == elem2 {
				isElement = true
				break
			}
		}

		if !isElement {
			count++
		}
	}
	return count
}
