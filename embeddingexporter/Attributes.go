package embeddingexporter

import (
	"encoding/json"
)

type Attributes map[string]any

func (a *Attributes) AsJson() ([]byte, error) {
	json, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}

	return json, nil
}

func (a *Attributes) AsString() (string, error) {
	json, err := a.AsJson()
	if err != nil {
		return "", err
	}

	return string(json), nil
}
