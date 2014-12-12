package multibayes

import (
	"encoding/json"
	"io/ioutil"
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

func LoadClassifierFromFile(filename string) (*Classifier, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return NewClassifierFromJSON(buf)
}

func (s *sparseColumn) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Data)
}

func (s *sparseColumn) UnmarshalJSON(buf []byte) error {
	var data []int

	err := json.Unmarshal(buf, &data)
	if err != nil {
		return err
	}

	s.Data = data

	return nil
}
