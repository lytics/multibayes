package multibayes

import (
	"encoding/json"
)

func (c *Classifier) MarshalJSON() ([]byte, error) {
	return json.Marshal(*c)
}

func NewClassifierFromJSON(buf []byte) (*Classifier, error) {
	tokenizer, _ := NewTokenizer(&TokenizerConf{
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
