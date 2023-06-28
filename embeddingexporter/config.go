package embeddingexporter

import (
	"fmt"

	"go.opentelemetry.io/collector/config/configtelemetry"
)

type Config struct {
	Verbosity   configtelemetry.Level `mapstructure:"verbosity,omitempty"`
	Embedding   EmbeddingConfig       `mapstructure:"embedding"`
	Persistence PersistenceConfig     `mapstructure:"persistence"`
}

type EmbeddingConfig struct {
	Endpoint     string            `mapstructure:"endpoint"`
	Key          string            `mapstructure:"key"`
	Version      string            `mapstructure:"version"`
	ModelMapping map[string]string `mapstructure:"model_mapping"`
}

type PersistenceConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Database int    `mapstructure:"database"`
}

func (cfg *Config) Validate() error {
	if cfg.Embedding.Endpoint == "" {
		return fmt.Errorf("embedding endpoint is required")
	}
	if cfg.Embedding.Key == "" {
		return fmt.Errorf("embedding key is required")
	}

	if cfg.Persistence.Host == "" {
		return fmt.Errorf("persistence host is required")
	}
	if cfg.Persistence.Port == "" {
		return fmt.Errorf("persistence port is required")
	}
	if cfg.Persistence.Password == "" {
		return fmt.Errorf("persistence password is required")
	}

	return nil
}
