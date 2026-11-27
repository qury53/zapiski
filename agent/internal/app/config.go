package app

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	GRPC_port    string        `envconfig:"GRPC_PORT" required:"true"`
	MetricsPort  string        `envconfig:"METRICS_PORT" default:"9092"`
	Env          string        `envconfig:"ENV" default:"dev"`
	ShutdownTime time.Duration `envconfig:"SHUTDOWN_TIME" default:"30s"`
}

func NewConfig() (*Config, error) {
	var config Config
	if err := envconfig.Process("AGENT", &config); err != nil {
		return nil, fmt.Errorf("failed process agent config: %w", err)
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
