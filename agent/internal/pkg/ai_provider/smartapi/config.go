package ai_provider_smartapi

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	SecretKey string `envconfig:"SECRET_KEY" required:"true"`
}

func NewConfig() (*Config, error) {
	var config Config
	if err := envconfig.Process("SMARTAPI", &config); err != nil {
		return nil, fmt.Errorf("smartapi ai provider config parse failed: %w", err)
	}

	return &config, nil
}

func NewConfigMust() *Config {
	config, err := NewConfig()
	if err != nil {
		panic(err)
	}

	return config
}
