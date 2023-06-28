package embeddingexporter

type Embeddings interface {
	Generate(input string) (Embedding, error)
}
