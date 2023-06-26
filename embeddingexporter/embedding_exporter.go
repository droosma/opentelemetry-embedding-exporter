package embeddingexporter

import (
	"context"

	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type embeddingExporter struct {
	// logsMarshaler    plog.Marshaler
	// metricsMarshaler pmetric.Marshaler
	// tracesMarshaler  ptrace.Marshaler
	embedding embedding
}

func newEmbeddingExporter() *embeddingExporter {
	//don't know how to get this from actual config
	config := NewConfig()

	return &embeddingExporter{
		// logsMarshaler:    otlptext.NewTextLogsMarshaler(),
		// metricsMarshaler: otlptext.NewTextMetricsMarshaler(),
		// tracesMarshaler:  otlptext.NewTextTracesMarshaler(),
		embedding: NewOpenAiEmbedder(config.OpenAiKey, config.OpenAiUri, config.OpenAiVersion),
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

func (s *embeddingExporter) pushLogs(_ context.Context, ld plog.Logs) error {
	//for now only extracting the text, the properties and stuff might also make sense.. maybe it doesnt.
	var input string
	rls := ld.ResourceLogs()
	for i := 0; i < rls.Len(); i++ {
		rl := rls.At(i)
		ills := rl.ScopeLogs()
		for j := 0; j < ills.Len(); j++ {
			ils := ills.At(j)
			logs := ils.LogRecords()
			for k := 0; k < logs.Len(); k++ {
				lr := logs.At(k)
				input += lr.Body().AsString()
			}
		}
	}

	//embedding, err := s.embedding.Embed(input)

	// if err != nil {
	// 	return err
	// }

	//TODO store the embeddings

	return nil
}
