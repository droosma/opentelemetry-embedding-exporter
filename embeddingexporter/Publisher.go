package embeddingexporter

type Publisher interface {
	Publish(embeddings []logEntryWithEmbedding) error
}
