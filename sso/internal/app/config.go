package application

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	GRPCPort     string        `envconfig:"GRPC_PORT" required:"true"`
	ShutdownTime time.Duration `envconfig:"SHUTDOWN_TIME" default:"30s"`
}

func NewConfig() (*Config, error) {
	var config Config
	if err := envconfig.Process("SSO", &config); err != nil {
		return nil, fmt.Errorf("process sso service config: %w", err)
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
