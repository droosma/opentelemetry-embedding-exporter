package embeddingexporter

import (
	"context"

	"go.opentelemetry.io/collector/pdata/pmetric"
)

type embeddingMetricsExporter struct {
	embedding   embedding
	persistence persistence
}

func newEmbeddingMetricsExporter(e embedding, p persistence) *embeddingMetricsExporter {
	return &embeddingMetricsExporter{
		embedding:   e,
		persistence: p,
	}
}

func (s *embeddingMetricsExporter) pushMetrics(_ context.Context, md pmetric.Metrics) error {
	//metrics might not make sense
	return nil
}
