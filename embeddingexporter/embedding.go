package embeddingexporter

import (
	"bytes"
	"encoding/binary"
)

type Embedding []*float32

func (e *Embedding) AsBytes() ([]byte, error) {
	values := make([]float32, len(*e))
	for i, v := range *e {
		values[i] = *v
	}

	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, values)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
