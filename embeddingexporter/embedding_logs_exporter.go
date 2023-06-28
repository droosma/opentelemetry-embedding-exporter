package embeddingexporter

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/collector/pdata/plog"
)

type embeddingLogsExporter struct {
	embedding   embedding
	persistence persistence
}

func newEmbeddingLogsExporter(e embedding, p persistence) *embeddingLogsExporter {
	return &embeddingLogsExporter{
		embedding:   e,
		persistence: p,
	}
}

func (s *embeddingLogsExporter) pushLogs(_ context.Context, ld plog.Logs) error {
	entries := extractLogEntries(ld)
	embeddings, embeddingsErrors := s.generateEmbeddingForLogEntries(entries)

	if embeddingsErrors != nil {
		return embeddingsErrors[0]
	}

	persistenceErrors := s.persistEmbeddings(embeddings)

	if persistenceErrors != nil {
		return persistenceErrors[0]
	}

	return nil
}

func (e logEntry) toLogEntryWithEmbedding(embedding []float32) logEntryWithEmbedding {
	return logEntryWithEmbedding{
		logEntry:  e,
		embedding: embedding,
	}
}

type logEntry struct {
	body       string
	level      string
	timestamp  time.Time
	TraceId    string
	SpanId     string
	attributes map[string]any
}

type logEntryWithEmbedding struct {
	logEntry  logEntry
	embedding []float32
}

func (s *embeddingLogsExporter) generateEmbeddingForLogEntries(entries []logEntry) ([]logEntryWithEmbedding, []error) {
	successesChan := make(chan logEntryWithEmbedding, len(entries))
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

			logEntry := entry.toLogEntryWithEmbedding(embedding)
			successesChan <- logEntry
		}(entry)
	}

	wg.Wait()

	close(successesChan)
	close(errorsChan)

	var successes []logEntryWithEmbedding
	for success := range successesChan {
		successes = append(successes, success)
	}

	var errors []error
	for err := range errorsChan {
		errors = append(errors, err)
	}

	return successes, errors
}

func (s *embeddingLogsExporter) persistEmbeddings(embeddings []logEntryWithEmbedding) []error {
	errorsChan := make(chan error, len(embeddings))

	var wg sync.WaitGroup

	for _, entry := range embeddings {
		wg.Add(1)

		go func(entry logEntryWithEmbedding) {
			defer wg.Done()

			key := fmt.Sprintf("log_%s_%s_%s",
				entry.logEntry.level,
				entry.logEntry.TraceId,
				uuid.New().String())

			properties := Properties{
				"timestamp": entry.logEntry.timestamp.Unix(),
				"body":      entry.logEntry.body,
				"level":     entry.logEntry.level,
				"traceId":   entry.logEntry.TraceId,
				"spanId":    entry.logEntry.SpanId,
			}
			err := properties.AddEmbedding(entry.embedding)
			if err != nil {
				errorsChan <- err
				return
			}
			err = properties.AddAttributes(entry.logEntry.attributes)
			if err != nil {
				errorsChan <- err
				return
			}

			err = s.persistence.Persist(key, properties)
			if err != nil {
				errorsChan <- err
				return
			}

		}(entry)
	}

	wg.Wait()

	close(errorsChan)

	var errors []error
	for err := range errorsChan {
		errors = append(errors, err)
	}

	return errors
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
					body:       lr.Body().AsString(),
					level:      lr.SeverityText(),
					timestamp:  lr.Timestamp().AsTime(),
					TraceId:    lr.TraceID().String(),
					SpanId:     lr.SpanID().String(),
					attributes: lr.Attributes().AsRaw(),
				})
			}
		}
	}
	return entries
}
