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
	embedding   Embeddings
	persistence Persistence
	publisher   Publisher
	mu          sync.Mutex
}

func (c *container) initialize(cfg *Config) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.embedding == nil && c.persistence == nil && c.publisher == nil {
		persistence := func(config PersistenceConfig) Persistence {
			return NewRedisPersistence(config.Host, config.Port, config.Password, config.Database)
		}

		embedding := func(config EmbeddingConfig) Embeddings {
			return NewAzureOpenAIEmbeddings(config.Key, config.Endpoint, config.ModelId, config.Version)
		}

		eventHub := func(config PublisherConfig) Publisher {
			if !config.Enabled {
				return disabledPublisher{}
			}
			return NewAzureEventHubPublisher(config.ConnectionString)
		}

		c.embedding = embedding(cfg.Embedding)
		c.persistence = persistence(cfg.Persistence)
		c.publisher = eventHub(cfg.Publisher)
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
		exporterhelper.WithTimeout(exporterhelper.TimeoutSettings{Timeout: 5}),
		exporterhelper.WithRetry(exporterhelper.RetrySettings{
			Enabled:         true,
			InitialInterval: 5,
			MaxInterval:     30,
		}),
		exporterhelper.WithQueue(exporterhelper.QueueSettings{
			Enabled:   true,
			QueueSize: 1000,
		}),
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
		exporterhelper.WithTimeout(exporterhelper.TimeoutSettings{Timeout: 5}),
		exporterhelper.WithRetry(exporterhelper.RetrySettings{
			Enabled:         true,
			InitialInterval: 5,
			MaxInterval:     30,
		}),
		exporterhelper.WithQueue(exporterhelper.QueueSettings{
			Enabled:   true,
			QueueSize: 1000,
		}),
	)
}

func (c *container) createLogsExporter(
	ctx context.Context,
	set exporter.CreateSettings,
	config component.Config) (exporter.Logs, error) {
	cfg := config.(*Config)
	c.initialize(cfg)
	x := newEmbeddingLogsExporter(c.embedding, c.persistence, c.publisher)
	return exporterhelper.NewLogsExporter(ctx, set, cfg,
		x.pushLogs,
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: false}),
		exporterhelper.WithTimeout(exporterhelper.TimeoutSettings{Timeout: 5}),
		exporterhelper.WithRetry(exporterhelper.RetrySettings{
			Enabled:         false,
			InitialInterval: 5,
			MaxInterval:     30,
		}),
		exporterhelper.WithQueue(exporterhelper.QueueSettings{
			Enabled:   false,
			QueueSize: 1000,
		}),
	)
}
