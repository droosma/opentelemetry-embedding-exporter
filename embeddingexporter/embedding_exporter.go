package embeddingexporter

import (
	"context"
	"sync"
	"time"

	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type embeddingExporter struct {
	embedding embedding
}

func newEmbeddingExporter(e embedding) *embeddingExporter {
	return &embeddingExporter{
		embedding: e,
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

type logEntry struct {
	body      string
	level     string
	timestamp time.Time
	TraceId   string
	SpanId    string
}

func (s *embeddingExporter) pushLogs(_ context.Context, ld plog.Logs) error {
	entries := extractLogEntries(ld)
	_, errors := s.processLogEntries(entries)

	if errors != nil {
		return errors[0]
	}

	return nil
}

func (s *embeddingExporter) processLogEntries(entries []logEntry) ([]float32, []error) {
	successesChan := make(chan []float32, len(entries))
	errorsChan := make(chan error, len(entries))

	var wg sync.WaitGroup

	for _, entry := range entries {
		wg.Add(1)

		go func(entry logEntry) {
			defer wg.Done()

			embedding, err := s.embedding.Embed(entry.body)
			if err != nil {
				errorsChan <- err
				return
			}

			successesChan <- embedding
		}(entry)
	}

	wg.Wait()

	close(successesChan)
	close(errorsChan)

	var successes []float32
	for success := range successesChan {
		successes = append(successes, success...)
	}

	var errors []error
	for err := range errorsChan {
		errors = append(errors, err)
	}

	return successes, errors
}

func extractLogEntries(ld plog.Logs) []logEntry {
	var entries []logEntry
	for i := 0; i < ld.ResourceLogs().Len(); i++ {
		rl := ld.ResourceLogs().At(i)
		for j := 0; j < rl.ScopeLogs().Len(); j++ {
			ils := rl.ScopeLogs().At(j)
			for k := 0; k < ils.LogRecords().Len(); k++ {
				lr := ils.LogRecords().At(k)
				entries = append(entries, logEntry{
					body:      lr.Body().AsString(),
					level:     lr.SeverityText(),
					timestamp: lr.Timestamp().AsTime(),
					TraceId:   lr.TraceID().String(),
					SpanId:    lr.SpanID().String(),
				})
			}
		}
	}
	return entries
}
