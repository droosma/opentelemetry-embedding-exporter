package embeddingexporter

import (
	"context"

	"go.opentelemetry.io/collector/pdata/ptrace"
)

type embeddingTracesExporter struct {
	embedding   embedding
	persistence persistence
}

func newEmbeddingTracesExporter(e embedding, p persistence) *embeddingTracesExporter {
	return &embeddingTracesExporter{
		embedding:   e,
		persistence: p,
	}
}

func (s *embeddingTracesExporter) pushTraces(_ context.Context, td ptrace.Traces) error {
	//traces might not make sense
	return nil
}
