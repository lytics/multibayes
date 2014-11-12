package forestclassifier

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/drewlanenga/multibayes/bag"
	"github.com/drewlanenga/multibayes/tokens"
	"github.com/lytics/lio/src/scores/mondrian"
	"github.com/ryanbressler/CloudForest"
)

// will have a forest for every class type
type ForestBag map[string]*mondrian.MondrianForest

// give it a training set? could be rewritten to be given a single line
func Learn(ngrams [][]tokens.NGram, classes [][]string) (ForestBag, map[string]*CloudForest.FeatureMatrix) {
	sparseMatrix := bag.NewSparseMatrix()
	for i, ngram := range ngrams {
		sparseMatrix.Add(ngram, classes[i])
	}

	matrices := sparseMatrix.ToFeatureMatrix()

	// forest bag has a forest for each document class
	forestBag := make(ForestBag)
	for class, classindex := range sparseMatrix.ClassMap {
		forestBag[class] = &mondrian.MondrianForest{
			Target: "C:" + class,
			Trees:  make(map[int64]*CloudForest.Tree),
		}

		newMap := make(map[string]int)
		for token, index := range matrices["tokens"].Map {
			newMap[token] = index
		}
		mapKey := "C:" + class
		newMap[mapKey] = len(matrices["tokens"].Map)

		newFeatures := make([]CloudForest.Feature, len(matrices["tokens"].Data)+1)
		for i, feature := range matrices["tokens"].Data {
			newFeatures[i] = feature
		}
		newFeatures[len(matrices["tokens"].Data)] = matrices["classes"].Data[classindex]

		featureMatrix := &CloudForest.FeatureMatrix{
			Data:       newFeatures,
			Map:        newMap,
			CaseLabels: make([]string, 0),
		}

		err := mondrian.GrowMondrianForest(featureMatrix, forestBag[class])
		if err != nil {
			fmt.Println(err)
		}
		mondrian.WriteForest(forestBag[class])

	}
	return forestBag, matrices
}

func (f ForestBag) Predict(tokenMatrix *CloudForest.FeatureMatrix, ngrams []tokens.NGram) map[string]float64 {
	newTokens := make(map[int]int)

	for _, ngram := range ngrams {
		gramString := ngram.String()
		gramString = "N:" + gramString

		if tokenIndex, ok := tokenMatrix.Map[gramString]; ok {
			newTokens[tokenIndex]++
		}
	}

	newTokenMap := make(map[string]int)
	for token, tokenindex := range tokenMatrix.Map {
		trimmed := strings.Trim(token, "N:")
		newTokenMap[trimmed] = tokenindex
	}

	newFeatures := bag.ToFeatures(newTokenMap, []map[int]int{newTokens}, "N:")
	featureMatrix := &CloudForest.FeatureMatrix{
		Data:       newFeatures,
		Map:        tokenMatrix.Map,
		CaseLabels: make([]string, 0),
	}

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
