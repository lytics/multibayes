Multibayes
==========

[![Build Status](https://travis-ci.org/lytics/multibayes.svg?branch=master)](https://travis-ci.org/lytics/multibayes) [![GoDoc](https://godoc.org/github.com/lytics/multibayes?status.svg)](https://godoc.org/github.com/lytics/multibayes)

Multiclass naive Bayesian document classification.

Often in document classification, a document may have more than one relevant classification -- a question on [stackoverflow](http://stackoverflow.com) might have tags "go", "map", and "interface".

While multinomial Bayesian classification offers a one-of-many classification, multibayes offers tools for many-of-many classification.  The multibayes library strives to offer efficient storage and calculation of multiple Bayesian posterior classification probabilities.

## Usage

A new classifier is created with the `NewClassifier` function, and can be trained by adding documents and classes by calling the `Add` method:

```go
classifier.Add("A new document", []string{"class1", "class2"})
```

Posterior probabilities for a new document are calculated by calling the `Posterior` method:

```go
classifier.Posterior("Another new document")
```

A posterior class probability is returned for each class observed in the training set, which the user can use to determine class assignment.  A user can then assign classifications according to his or her own heuristics -- for example, by using all classes that yield a posterior probability greater than 0.8


## Example

```go
documents := []struct {
	Text    string
	Classes []string
}{
	{
		Text:    "My dog has fleas.",
		Classes: []string{"vet"},
	},
	{
		Text:    "My cat has ebola.",
		Classes: []string{"vet", "cdc"},
	},
	{
		Text:    "Aaron has ebola.",
		Classes: []string{"cdc"},
	},
}

classifier := NewClassifier()
classifier.MinClassSize = 0

// train the classifier
for _, document := range documents {
	classifier.Add(document.Text, document.Classes)
}

// predict new classes
probs := classifier.Posterior("Aaron's dog has fleas.")
fmt.Printf("Posterior Probabilities: %+v\n", probs)

// Posterior Probabilities: map[vet:0.8571 cdc:0.2727]
```
