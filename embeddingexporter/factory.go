package embeddingexporter

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

const (
	typeStr = "embeddingexporter"
)

func NewFactory() exporter.Factory {
	return exporter.NewFactory(
		typeStr,
		createDefaultConfig,
		exporter.WithTraces(createTracesExporter, component.StabilityLevelDevelopment),
		exporter.WithMetrics(createMetricsExporter, component.StabilityLevelDevelopment),
		exporter.WithLogs(createLogsExporter, component.StabilityLevelDevelopment),
	)
}

func createDefaultConfig() component.Config {
	return NewConfig()
}

func createTracesExporter(ctx context.Context, set exporter.CreateSettings, config component.Config) (exporter.Traces, error) {
	cfg := config.(*Config)
	s := newEmbeddingExporter()
	return exporterhelper.NewTracesExporter(ctx, set, cfg,
		s.pushTraces,
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: false}),
		exporterhelper.WithTimeout(exporterhelper.TimeoutSettings{Timeout: 0}),
		exporterhelper.WithRetry(exporterhelper.RetrySettings{Enabled: false}),
		exporterhelper.WithQueue(exporterhelper.QueueSettings{Enabled: false}),
	)
}

func createMetricsExporter(ctx context.Context, set exporter.CreateSettings, config component.Config) (exporter.Metrics, error) {
	cfg := config.(*Config)
	s := newEmbeddingExporter()
	return exporterhelper.NewMetricsExporter(ctx, set, cfg,
		s.pushMetrics,
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: false}),
		exporterhelper.WithTimeout(exporterhelper.TimeoutSettings{Timeout: 0}),
		exporterhelper.WithRetry(exporterhelper.RetrySettings{Enabled: false}),
		exporterhelper.WithQueue(exporterhelper.QueueSettings{Enabled: false}),
	)
}

func createLogsExporter(ctx context.Context, set exporter.CreateSettings, config component.Config) (exporter.Logs, error) {
	cfg := config.(*Config)
	s := newEmbeddingExporter()
	return exporterhelper.NewLogsExporter(ctx, set, cfg,
		s.pushLogs,
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: false}),
		exporterhelper.WithTimeout(exporterhelper.TimeoutSettings{Timeout: 0}),
		exporterhelper.WithRetry(exporterhelper.RetrySettings{Enabled: false}),
		exporterhelper.WithQueue(exporterhelper.QueueSettings{Enabled: false}),
	)
}
