package embeddingexporter

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

type embeddingMetricsExporter struct {
	embedding   Embeddings
	persistence Persistence
}

type metricEntry struct {
	name        string
	description string
	unit        string
	dataType    string
}

type metricEntryWithEmbedding struct {
	metricEntry metricEntry
	embedding   []float32
}

func (e metricEntry) toMetricEntryWithEmbedding(embedding []float32) metricEntryWithEmbedding {
	return metricEntryWithEmbedding{
		metricEntry: e,
		embedding:   embedding,
	}
}

func newEmbeddingMetricsExporter(e Embeddings, p Persistence) *embeddingMetricsExporter {
	return &embeddingMetricsExporter{
		embedding:   e,
		persistence: p,
	}
}

func (s *embeddingMetricsExporter) pushMetrics(_ context.Context, md pmetric.Metrics) error {
	entries := extractMetricEntries(md)
	errors := s.persistEntries(entries)

	if errors != nil {
		return errors[0]
	}

	return nil
}

func extractMetricEntries(md pmetric.Metrics) []metricEntry {
	var entries []metricEntry
	for i := 0; i < md.ResourceMetrics().Len(); i++ {
		rl := md.ResourceMetrics().At(i)
		for j := 0; j < rl.ScopeMetrics().Len(); j++ {
			ils := rl.ScopeMetrics().At(j)
			for k := 0; k < ils.Metrics().Len(); k++ {
				lr := ils.Metrics().At(k)
				/*
					Not sure what to do with metric values, for reference see:
					https://github.com/open-telemetry/opentelemetry-collector/blob/main/exporter/loggingexporter/internal/otlptext/databuffer.go#L66C30-L66C30

					metricType := lr.Type()
					switch metricType {
					case pmetric.MetricTypeEmpty:
						break
					case pmetric.MetricTypeGauge:
						//lr.Gauge().DataPoints().At(0).Value()
						break
					case pmetric.MetricTypeSum:
						//lr.Sum().DataPoints().At(0).Value()
						break
					case pmetric.MetricTypeHistogram:
						//lr.Histogram().DataPoints().At(0).Value()
						break
					case pmetric.MetricTypeExponentialHistogram:
						//lr.ExponentialHistogram().DataPoints().At(0).Value()
						break
					case pmetric.MetricTypeSummary:
						//lr.Summary().DataPoints().At(0).Value()
						break
					}
				*/

				entries = append(entries, metricEntry{
					name:        lr.Name(),
					description: lr.Description(),
					unit:        lr.Unit(),
					dataType:    lr.Type().String(),
				})
			}
		}
	}
	return entries
}

func (s *embeddingMetricsExporter) persistEntries(entries []metricEntry) []error {
	errorsChan := make(chan error, len(entries))

	var wg sync.WaitGroup

	for _, entry := range entries {
		wg.Add(1)

		go func(entry metricEntry) {
			defer wg.Done()

			key := fmt.Sprintf("metric_%s_%s",
				entry.name,
				uuid.New().String())

			properties := Properties{
				"name":        entry.name,
				"description": entry.description,
				"unit":        entry.unit,
				"dataType":    entry.dataType,
			}

			err := s.persistence.Persist(key, properties)
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
