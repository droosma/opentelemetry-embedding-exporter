package embeddingexporter

import (
	"context"

	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type embeddingExporter struct {
	logsMarshaler    plog.Marshaler
	metricsMarshaler pmetric.Marshaler
	tracesMarshaler  ptrace.Marshaler
}

func newEmbeddingExporter() {

}

func (s *embeddingExporter) pushTraces(_ context.Context, td ptrace.Traces) error {
	buf, err := s.tracesMarshaler.MarshalTraces(td)
	if err != nil {
		return err
	}
	return nil
}

func (s *embeddingExporter) pushMetrics(_ context.Context, md pmetric.Metrics) error {
	buf, err := s.metricsMarshaler.MarshalMetrics(md)
	if err != nil {
		return err
	}
	return nil
}

func (s *embeddingExporter) pushLogs(_ context.Context, ld plog.Logs) error {
	buf, err := s.logsMarshaler.MarshalLogs(ld)
	if err != nil {
		return err
	}
	return nil
}
