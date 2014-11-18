package multibayes

import (
	"encoding/json"
)

type jsonableClassifier struct {
	Matrix *sparseMatrix `json:"matrix"`
}

func (c *Classifier) MarshalJSON() ([]byte, error) {
	return json.Marshal(&jsonableClassifier{c.Matrix})
}

func (c *Classifier) UnmarshalJSON(buf []byte) error {
	j := jsonableClassifier{}

	err := json.Unmarshal(buf, &j)
	if err != nil {
		return nil
	}

	*c = *NewClassifier()
	c.Matrix = j.Matrix

	return nil
}

// Initialize a new classifier from a JSON byte slice.
func NewClassifierFromJSON(buf []byte) (*Classifier, error) {
	classifier := &Classifier{}

	err := classifier.UnmarshalJSON(buf)
	if err != nil {
		return nil, err
	}

	return classifier, nil
}
