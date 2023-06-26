package embeddingexporter

import (
	"fmt"

	"go.opentelemetry.io/collector/config/configtelemetry"
)

type Config struct {
	Verbosity configtelemetry.Level `mapstructure:"verbosity,omitempty"`
	Embedding EmbeddingConfig       `mapstructure:"embedding"`
}

type EmbeddingConfig struct {
	Endpoint string `mapstructure:"endpoint"`
	Key      string `mapstructure:"key"`
	Version  string `mapstructure:"version"`
}

func NewConfig() Config {
	return Config{
		Verbosity: configtelemetry.LevelNormal,
		Embedding: EmbeddingConfig{
			Version:  "2023-05-15",
			Key:      "5c8e2b1c28414c5183185154c66c1242",
			Endpoint: "https://rg-openai-sandbox.openai.azure.com/",
		},
	}
}

func (cfg *Config) Validate() error {
	if cfg.Embedding.Endpoint == "" {
		return fmt.Errorf("embedding endpoint is required")
	}
	if cfg.Embedding.Key == "" {
		return fmt.Errorf("embedding key is required")
	}
	return nil
}
