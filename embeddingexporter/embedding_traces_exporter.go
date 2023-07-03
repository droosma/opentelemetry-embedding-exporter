package embeddingexporter

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"go.opentelemetry.io/collector/pdata/ptrace"
)

type embeddingTracesExporter struct {
	embedding   Embeddings
	persistence Persistence
}

func newEmbeddingTracesExporter(e Embeddings, p Persistence) *embeddingTracesExporter {
	return &embeddingTracesExporter{
		embedding:   e,
		persistence: p,
	}
}

func (s *embeddingTracesExporter) pushTraces(_ context.Context, td ptrace.Traces) error {
	entries := extractTraceEntries(td)
	embeddings, errors := s.generateEmbeddingForTraceEntries(entries)

	if errors != nil {
		return errors[0]
	}

	errors = s.persistEmbeddings(embeddings)

	if errors != nil {
		return errors[0]
	}

	return nil
}

type traceEntry struct {
	id         string
	name       string
	kind       string
	start      time.Time
	end        time.Time
	attributes Attributes
	status     string
	message    string
	TraceId    string
	SpanId     string
}

func (e traceEntry) toTraceEntryWithEmbedding(embedding Embedding) traceEntryWithEmbedding {
	return traceEntryWithEmbedding{
		traceEntry: e,
		embedding:  embedding,
	}
}

func (e traceEntry) embeddingBody() string {
	var builder strings.Builder

	builder.WriteString(e.name + " ")
	builder.WriteString(e.kind + " ")
	builder.WriteString(e.status + " ")
	builder.WriteString(e.message + " ")

	attrString, err := e.attributes.AsString()
	if err == nil {
		builder.WriteString(attrString)
	}

	return builder.String()
}

type traceEntryWithEmbedding struct {
	traceEntry traceEntry
	embedding  Embedding
}

func (s *embeddingTracesExporter) generateEmbeddingForTraceEntries(entries []traceEntry) ([]traceEntryWithEmbedding, []error) {
	successesChan := make(chan traceEntryWithEmbedding, len(entries))
	errorsChan := make(chan error, len(entries))

	var wg sync.WaitGroup

	for _, entry := range entries {
		wg.Add(1)

		go func(entry traceEntry) {
			defer wg.Done()

			embedding, err := s.embedding.Generate(entry.embeddingBody())
			if err != nil {
				errorsChan <- err
				return
			}

			logEntry := entry.toTraceEntryWithEmbedding(embedding)
			successesChan <- logEntry
		}(entry)
	}

	wg.Wait()

	close(successesChan)
	close(errorsChan)

	var successes []traceEntryWithEmbedding
	for success := range successesChan {
		successes = append(successes, success)
	}

	var errors []error
	for err := range errorsChan {
		errors = append(errors, err)
	}

	return successes, errors
}

func extractTraceEntries(td ptrace.Traces) []traceEntry {
	var entries []traceEntry
	for i := 0; i < td.ResourceSpans().Len(); i++ {
		rl := td.ResourceSpans().At(i)
		for j := 0; j < rl.ScopeSpans().Len(); j++ {
			ils := rl.ScopeSpans().At(j)
			for k := 0; k < ils.Spans().Len(); k++ {
				lr := ils.Spans().At(k)
				entries = append(entries, traceEntry{
					TraceId:    lr.TraceID().String(),
					SpanId:     lr.SpanID().String(),
					start:      lr.StartTimestamp().AsTime(),
					end:        lr.EndTimestamp().AsTime(),
					status:     lr.Status().Code().String(),
					message:    lr.Status().Message(),
					id:         lr.SpanID().String(),
					name:       lr.Name(),
					kind:       lr.Kind().String(),
					attributes: lr.Attributes().AsRaw(),
				})
			}
		}
	}
	return entries
}

func (s *embeddingTracesExporter) persistEmbeddings(embeddings []traceEntryWithEmbedding) []error {
	errorsChan := make(chan error, len(embeddings))

	var wg sync.WaitGroup

	for _, entry := range embeddings {
		wg.Add(1)

		go func(entry traceEntryWithEmbedding) {
			defer wg.Done()
			key := fmt.Sprintf("trace_%s_%s", entry.traceEntry.status, entry.traceEntry.id)

			properties := Properties{
				"id":      entry.traceEntry.id,
				"name":    entry.traceEntry.name,
				"kind":    entry.traceEntry.kind,
				"start":   entry.traceEntry.start.Unix(),
				"end":     entry.traceEntry.end.Unix(),
				"status":  entry.traceEntry.status,
				"message": entry.traceEntry.message,
				"traceId": entry.traceEntry.TraceId,
				"spanId":  entry.traceEntry.SpanId,
			}

			err := properties.AddEmbedding(entry.embedding)
			if err != nil {
				errorsChan <- err
				return
			}
			err = properties.AddAttributes(entry.traceEntry.attributes)
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
