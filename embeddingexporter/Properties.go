package embeddingexporter

import (
	"encoding/json"
)

type Properties map[string]interface{}

func (p *Properties) AddEmbedding(embedding Embedding) error {
	bytes, error := embedding.AsBytes()

	if error != nil {
		return error
	}

	(*p)["embedding"] = bytes

	return nil
}

func (p *Properties) AddAttributes(attributes map[string]any) error {
	attributesJson, err := json.Marshal(attributes)
	if err != nil {
		return err
	}

	(*p)["attributes"] = attributesJson

	return nil
}
