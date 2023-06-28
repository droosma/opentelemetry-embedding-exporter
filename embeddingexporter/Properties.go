package embeddingexporter

type Properties map[string]interface{}

func (p *Properties) AddEmbedding(embedding Embedding) error {
	bytes, error := embedding.AsBytes()

	if error != nil {
		return error
	}

	(*p)["embedding"] = bytes

	return nil
}

func (p *Properties) AddAttributes(attributes Attributes) error {
	attributesJson, err := attributes.AsJson()
	if err != nil {
		return err
	}

	(*p)["attributes"] = attributesJson

	return nil
}
