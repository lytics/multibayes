package multibayes

import (
	"encoding/json"
)

// Serialize the classifier to JSON.
func (c *Classifier) MarshalJSON() ([]byte, error) {
	return json.Marshal(*c)
}

// Initialize a new classifier from a JSON byte slice.
func NewClassifierFromJSON(buf []byte) (*Classifier, error) {
	tokenizer, _ := newTokenizer(&tokenizerConf{
		NGramSize: 1,
	})

	classifier := &Classifier{
		Tokenizer: tokenizer,
	}

	err := json.Unmarshal(buf, classifier)
	if err != nil {
		return nil, err
	}

	return classifier, nil
}
