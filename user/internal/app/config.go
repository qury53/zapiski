package app

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	GRPCPort     string        `envconfig:"GRPC_PORT" required:"true"`
	MetricsPort  string        `envconfig:"METRICS_PORT" default:"9093"`
	ShutdownTime time.Duration `envconfig:"SHUTDOWN_TIME" default:"30s"`
	Env          string        `envconfig:"ENV" default:"dev"`
}

// NewConfig создает конфиг или бан
//
// Returns: *Config error

func NewConfig() (*Config, error) {
	var config Config
	if err := envconfig.Process("USER", &config); err != nil {
		return nil, fmt.Errorf("failed to process user config: %w", err)
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
