package embeddingexporter

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
)

type Properties map[string]interface{}

func (p *Properties) AddEmbedding(embedding []float32) error {
	asBytes := func(floats []float32) ([]byte, error) {
		buf := new(bytes.Buffer)
		err := binary.Write(buf, binary.LittleEndian, floats)
		if err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}

	bytes, error := asBytes(embedding)

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
