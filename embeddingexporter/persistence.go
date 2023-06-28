package embeddingexporter

type Persistence interface {
	Persist(key string, properties Properties) error
}
