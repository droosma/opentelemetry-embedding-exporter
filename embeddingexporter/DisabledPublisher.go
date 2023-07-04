package embeddingexporter

import (
	"fmt"
)

type disabledPublisher struct{}

func (np disabledPublisher) Publish(embeddings []logEntryWithEmbedding) error {
	fmt.Println("Publisher is disabled")
	return nil
}
