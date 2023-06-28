package embeddingexporter

import (
	"bytes"
	"encoding/binary"
)

type Embedding []float32

func (e *Embedding) AsBytes() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, e)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
