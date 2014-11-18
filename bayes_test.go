package multibayes

import (
	"testing"

	//"github.com/bmizerany/assert"
	"github.com/drewlanenga/multibayes/testutil"
)

func TestPosterior(t *testing.T) {
	classifier := NewClassifier()

	testdata := testutil.GetTestData()
	for i, _ := range testdata.Docs {
		classifier.Add(testdata.Docs[i], testdata.Classes[i])
	}

	for i, _ := range testdata.Docs {
		_ = classifier.Posterior(testdata.Docs[i])
	}
}
