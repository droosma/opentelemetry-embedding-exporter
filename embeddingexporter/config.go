package embeddingexporter

import (
	"fmt"

	"go.opentelemetry.io/collector/config/configtelemetry"
)

type Config struct {
	Verbosity   configtelemetry.Level `mapstructure:"verbosity,omitempty"`
	Embedding   EmbeddingConfig       `mapstructure:"embedding"`
	Persistence PersistenceConfig     `mapstructure:"persistence"`
	Publisher   PublisherConfig       `mapstructure:"publisher, omitempty"`
}

type EmbeddingConfig struct {
	Endpoint string `mapstructure:"endpoint"`
	Key      string `mapstructure:"key"`
	Version  string `mapstructure:"version"`
	ModelId  string `mapstructure:"model_id"`
}

type PersistenceConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Database int    `mapstructure:"database"`
}

type PublisherConfig struct {
	Enabled          bool   `mapstructure:"enabled"`
	ConnectionString string `mapstructure:"connection_string"`
}

func (cfg *Config) Validate() error {
	if cfg.Embedding.Endpoint == "" {
		return fmt.Errorf("embedding endpoint is required")
	}
	if cfg.Embedding.Key == "" {
		return fmt.Errorf("embedding key is required")
	}
	if cfg.Embedding.ModelId == "" {
		return fmt.Errorf("embedding model_id is required")
	}

	if cfg.Persistence.Host == "" {
		return fmt.Errorf("persistence host is required")
	}
	if cfg.Persistence.Port == "" {
		return fmt.Errorf("persistence port is required")
	}

	if cfg.Publisher.Enabled && cfg.Publisher.ConnectionString == "" {
		return fmt.Errorf("publisher connection string is required")
	}

	return nil
}
