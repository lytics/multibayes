package multibayes

import (
	"math"
)

var (
	smoother = 1 // laplace
)

type Classifier struct {
	Tokenizer *Tokenizer    `json:"-"`
	Matrix    *SparseMatrix `json:"matrix"`
}

func NewClassifier() *Classifier {
	tokenizer, _ := NewTokenizer(&TokenizerConf{
		NGramSize: 1,
	})

	sparse := NewSparseMatrix()

	return &Classifier{
		Tokenizer: tokenizer,
		Matrix:    sparse,
	}
}

func (c *Classifier) Add(document string, classes []string) {
	ngrams := c.Tokenizer.Parse(document)
	c.Matrix.Add(ngrams, classes)
}

func (c *Classifier) Posterior(document string) map[string]float64 {
	tokens := c.Tokenizer.Parse(document)
	predictions := make(map[string]float64)

	for class, classcolumn := range c.Matrix.Classes {
		n := classcolumn.Count()

		priors := []float64{
			float64(n+smoother) / float64(c.Matrix.N+(smoother*2)),            // P(C=Y)
			float64(c.Matrix.N-n+smoother) / float64(c.Matrix.N+(smoother*2)), // P(C=N)
		}

		loglikelihood := []float64{1.0, 1.0}

		// check if each token is in our token sparse matrix
		for _, token := range tokens {
			if tokencolumn, ok := c.Matrix.Tokens[token.String()]; ok {
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
