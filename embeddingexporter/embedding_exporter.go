package embeddingexporter

import (
	"context"

	"go.opentelemetry.io/collector/pdata/pmetric"
)

type embeddingExporter struct {
	embedding   embedding
	persistence persistence
}

func newEmbeddingExporter(e embedding, p persistence) *embeddingExporter {
	return &embeddingExporter{
		embedding:   e,
		persistence: p,
	}
}

func (s *embeddingExporter) pushMetrics(_ context.Context, md pmetric.Metrics) error {
	//metrics might not make sense
	return nil
}
