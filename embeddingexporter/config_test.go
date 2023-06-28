package embeddingexporter

import (
	"testing"

	"go.opentelemetry.io/collector/config/configtelemetry"
)

func validConfig() *Config {
	return &Config{
		Verbosity: configtelemetry.LevelDetailed,
		Embedding: EmbeddingConfig{
			Endpoint: "https://example.com",
			Key:      "key",
			Version:  "version",
			ModelMapping: map[string]string{
				"gpt-3.5-turbo":          "turbo",
				"text-embedding-ada-002": "embedding",
			},
		},
		Persistence: PersistenceConfig{
			Host:     "localhost",
			Port:     "6379",
			Password: "password",
			Database: 0,
		},
	}
}

func TestConfig_Validate_ValidConfig(t *testing.T) {
	cfg := validConfig()
	err := cfg.Validate()

	if err != nil {
		t.Errorf("Expected not error for valid Config")
	}
}

func TestConfig_Validate_Embedding_Endpoint(t *testing.T) {
	cfg := validConfig()
	cfg.Embedding.Endpoint = ""
	err := cfg.Validate()

	if err == nil {
		t.Errorf("Expected error for missing Embedding.Endpoint, got nil")
	}
}

func TestConfig_Validate_Embedding_Key(t *testing.T) {
	cfg := validConfig()
	cfg.Embedding.Key = ""

	err := cfg.Validate()
	if err == nil {
		t.Errorf("Expected error for missing Embedding.Key, got nil")
	}
}

func TestConfig_Validate_Persistence_Host(t *testing.T) {
	cfg := validConfig()
	cfg.Persistence.Host = ""

	err := cfg.Validate()
	if err == nil {
		t.Errorf("Expected error for missing Persistence.Host, got nil")
	}
}

func TestConfig_Validate_Persistence_Port(t *testing.T) {
	cfg := validConfig()
	cfg.Persistence.Port = ""

	err := cfg.Validate()
	if err == nil {
		t.Errorf("Expected error for missing Persistence.Port, got nil")
	}
}
