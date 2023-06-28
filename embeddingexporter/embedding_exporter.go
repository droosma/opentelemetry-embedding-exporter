package embeddingexporter

import (
	"context"

	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
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

func (s *embeddingExporter) pushTraces(_ context.Context, td ptrace.Traces) error {
	//traces might not make sense
	return nil
}

func (s *embeddingExporter) pushMetrics(_ context.Context, md pmetric.Metrics) error {
	//metrics might not make sense
	return nil
}
