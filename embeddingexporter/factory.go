package embeddingexporter

import (
	"context"
	"sync"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configtelemetry"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

const (
	typeStr = "embeddingexporter"
)

type container struct {
	embedding   embedding
	persistence persistence
	mu          sync.Mutex
}

func (c *container) initialize(cfg *Config) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.embedding == nil && c.persistence == nil {
		persistence := func(config PersistenceConfig) persistence {
			return NewRedisPersistence(config.Host, config.Port, config.Password, config.Database)
		}

		embedding := func(config EmbeddingConfig) embedding {
			return NewOpenAiEmbedder(config.Key, config.Endpoint, config.ModelMapping, config.Version)
		}

		c.embedding = embedding(cfg.Embedding)
		c.persistence = persistence(cfg.Persistence)
	}
}

func NewFactory() exporter.Factory {
	c := &container{}

	return exporter.NewFactory(
		typeStr,
		createDefaultConfig,
		exporter.WithTraces(c.createTracesExporter, component.StabilityLevelDevelopment),
		exporter.WithMetrics(c.createMetricsExporter, component.StabilityLevelDevelopment),
		exporter.WithLogs(c.createLogsExporter, component.StabilityLevelDevelopment),
	)
}

func createDefaultConfig() component.Config {
	return &Config{
		Verbosity: configtelemetry.LevelNormal,
	}
}

func (c *container) createTracesExporter(
	ctx context.Context,
	set exporter.CreateSettings,
	config component.Config) (exporter.Traces, error) {
	cfg := config.(*Config)
	c.initialize(cfg)
	x := newEmbeddingTracesExporter(c.embedding, c.persistence)
	return exporterhelper.NewTracesExporter(ctx, set, cfg,
		x.pushTraces,
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: false}),
		exporterhelper.WithTimeout(exporterhelper.TimeoutSettings{Timeout: 0}),
		exporterhelper.WithRetry(exporterhelper.RetrySettings{Enabled: false}),
		exporterhelper.WithQueue(exporterhelper.QueueSettings{Enabled: false}),
	)
}

func (c *container) createMetricsExporter(
	ctx context.Context,
	set exporter.CreateSettings,
	config component.Config) (exporter.Metrics, error) {
	cfg := config.(*Config)
	c.initialize(cfg)
	x := newEmbeddingMetricsExporter(c.embedding, c.persistence)
	return exporterhelper.NewMetricsExporter(ctx, set, cfg,
		x.pushMetrics,
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: false}),
		exporterhelper.WithTimeout(exporterhelper.TimeoutSettings{Timeout: 0}),
		exporterhelper.WithRetry(exporterhelper.RetrySettings{Enabled: false}),
		exporterhelper.WithQueue(exporterhelper.QueueSettings{Enabled: false}),
	)
}

func (c *container) createLogsExporter(
	ctx context.Context,
	set exporter.CreateSettings,
	config component.Config) (exporter.Logs, error) {
	cfg := config.(*Config)
	c.initialize(cfg)
	x := newEmbeddingLogsExporter(c.embedding, c.persistence)
	return exporterhelper.NewLogsExporter(ctx, set, cfg,
		x.pushLogs,
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: false}),
		exporterhelper.WithTimeout(exporterhelper.TimeoutSettings{Timeout: 0}),
		exporterhelper.WithRetry(exporterhelper.RetrySettings{Enabled: false}),
		exporterhelper.WithQueue(exporterhelper.QueueSettings{Enabled: false}),
	)
}
