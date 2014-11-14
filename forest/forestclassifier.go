package forestclassifier

import (
	"fmt"
	"strconv"
	//"strings"

	"github.com/drewlanenga/multibayes/matrix"
	"github.com/drewlanenga/multibayes/tokens"
	"github.com/lytics/lio/src/scores/mondrian"
	"github.com/ryanbressler/CloudForest"
)

// fix so that every token header has the prefix "N:" (and every class "C:")
// before building the forest and using it for predictions

// will have a forest for every class type
type ForestBag map[string]*mondrian.MondrianForest

// give it a training set? could be rewritten to be given a single line
func Learn(ngrams [][]tokens.NGram, classes [][]string) (ForestBag, map[string]*CloudForest.FeatureMatrix) {
	sparseMatrix := matrix.NewSparseMatrix()
	for i, ngram := range ngrams {
		sparseMatrix.Add(ngram, classes[i])
	}

	matrices := sparseMatrix.ToFeatureMatrix()
	// forest bag has a forest for each document class
	forestBag := make(ForestBag)
	for class, _ := range sparseMatrix.Classes {
		forestBag[class] = &mondrian.MondrianForest{
			Target: class,
			Trees:  make(map[int64]*CloudForest.Tree),
		}

		newMap := make(map[string]int)
		for token, index := range matrices["tokens"].Map {
			newMap[token] = index
		}
		mapKey := class
		newMap[mapKey] = len(matrices["tokens"].Map)

		newFeatures := make([]CloudForest.Feature, len(matrices["tokens"].Data)+1)
		for i, feature := range matrices["tokens"].Data {
			newFeatures[i] = feature
		}
		newFeatures[len(matrices["tokens"].Data)] = sparseMatrix.Classes[class].ToFeature(sparseMatrix.N)

		featureMatrix := &CloudForest.FeatureMatrix{
			Data:       newFeatures,
			Map:        newMap,
			CaseLabels: make([]string, 0),
		}

		err := mondrian.GrowMondrianForest(featureMatrix, forestBag[class])
		if err != nil {
			fmt.Println(err)
		}
		//mondrian.WriteForest(forestBag[class])

	}
	return forestBag, matrices
}

func (f ForestBag) Predict(tokenMatrix *CloudForest.FeatureMatrix, ngrams []tokens.NGram) map[string]float64 {
	newTokens := make(map[int]int)

	for _, ngram := range ngrams {
		gramString := ngram.String()
		//gramString = "N:" + gramString

		if tokenIndex, ok := tokenMatrix.Map[gramString]; ok {
			newTokens[tokenIndex]++
		}
	}

	newTokenMap := make(map[string]int)
	for token, tokenindex := range tokenMatrix.Map {
		newTokenMap[token] = tokenindex
	}

	// create a single row feature matrix with
	// all the columns of the original feature matrix
	Tokens := []map[int]int{newTokens}
	newFeatures := make([]CloudForest.Feature, len(newTokenMap))
	tokenCount := make([]float64, len(Tokens))
	for tokenname, tokenindex := range newTokenMap {
		for rowindex, columnmap := range Tokens {
			if _, ok := columnmap[tokenindex]; !ok {
				tokenCount[rowindex] = 0
			} else {
				tokenCount[rowindex] = float64(columnmap[tokenindex])
			}
		}
		f := &CloudForest.DenseNumFeature{
			NumData:    tokenCount,
			Missing:    make([]bool, len(Tokens)),
			Name:       tokenname,
			HasMissing: false,
		}
		newFeatures[tokenindex] = f
	}
	///

	featureMatrix := &CloudForest.FeatureMatrix{
		Data:       newFeatures,
		Map:        tokenMatrix.Map,
		CaseLabels: make([]string, 0),
	}
	//fmt.Println(featureMatrix.Map)
	//fmt.Println(tokenMatrix.Map)

	predictions := make(map[string]float64)
	for class, forest := range f {

		var bb CloudForest.VoteTallyer
		bb = CloudForest.NewNumBallotBox(featureMatrix.Data[0].Length())

		for _, tree := range forest.Trees {
			tree.Vote(featureMatrix, bb)
		}

		targeti, hasTarget := featureMatrix.Map[forest.Target]
		if hasTarget {
			// note this is numerical error, not a functional error
			er := bb.TallyError(featureMatrix.Data[targeti])
			if er == float64(0) {
				// this would indicate a (perfectly) overfit model
				// not likely, not sure yet how we'd mitigate against it
			}
		}

		pred, _ := strconv.ParseFloat(bb.Tally(0), 64)
		predictions[class] = pred
	}
	return predictions
}
