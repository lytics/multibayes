package leafytree

import (
	//"fmt"

	"github.com/drewlanenga/multibayes/tokens"
)

var (
	// laplacian smoother by default
	smoother = 1
)

type Leaf struct {
	Classes map[string]int
	Count   int
}

func NewLeaf() *Leaf {
	return &Leaf{
		Classes: make(map[string]int),
	}
}

func (l *Leaf) Learn(classes []string) {
	for _, class := range classes {
		l.Classes[class]++
	}
	l.Count++
}

// Document Gram Matrix
type Tree struct {
	Leaves        map[string]*Leaf // leaves are features/tokens which contain classes
	TotalFeatures *Leaf
	TotalClasses  *Leaf
}

func NewTree() *Tree {
	return &Tree{
		Leaves:        make(map[string]*Leaf),
		TotalFeatures: NewLeaf(),
		TotalClasses:  NewLeaf(),
	}
}

// []tokens.Ngram is essentially a document
func (t *Tree) Learn(ngrams []tokens.NGram, classes []string) {
	featureStrings := ngramsToStrings(ngrams)

	// update each class leaf
	for _, class := range classes {
		leaf, ok := t.Leaves[class]
		if !ok {
			leaf = NewLeaf()
		}

		leaf.Learn(featureStrings)
		t.Leaves[class] = leaf
	}

	t.TotalClasses.Learn(classes)
	t.TotalFeatures.Learn(featureStrings)
}

func (t *Tree) Predict(ngrams []tokens.NGram) map[string]float64 {
	// P(C|F) = ( P(F|C) * P(C) ) / P(F)
	featureStrings := ngramsToStrings(ngrams)

	// predict a vector of class probabilites
	probF := 0.0
	probCF := make(map[string]float64)
	for class, classCount := range t.TotalClasses.Classes {
		probC := float64(classCount) / float64(t.TotalClasses.Count)

		probFC := probablize(t.Leaves[class].Classes, featureStrings, t.Leaves[class].Count)

		probCF[class] = (probFC * probC)
		probF += probCF[class]
	}

	// scale by the sum of the probs
	for class, classProb := range probCF {
		probCF[class] = classProb / probF
	}

	return probCF
}

func probablize(haystack map[string]int, needles []string, n int) float64 {
	extracted := extract(haystack, needles)
	scaled := scale(extracted, n)

	return product(scaled)
}

func extract(haystack map[string]int, needles []string) []int {
	result := make([]int, len(needles))
	for i, needle := range needles {
		value, ok := haystack[needle]
		if ok {
			result[i] = value
		}
	}
	return result
}

func scale(vector []int, n int) []float64 {
	result := make([]float64, len(vector))
	for i, v := range vector {
		result[i] = float64(v+smoother) / float64(n+smoother)
	}
	return result
}

func product(vector []float64) float64 {
	value := 1.0
	for _, v := range vector {
		value *= v
	}
	return value
}

func ngramsToStrings(ngrams []tokens.NGram) []string {
	ngramStrings := make([]string, len(ngrams))
	for i, ngram := range ngrams {
		ngramStrings[i] = ngram.String()
	}
	return ngramStrings
}
