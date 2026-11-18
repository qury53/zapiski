package redis

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Addr       string        `envconfig:"ADDR"`
	Password   string        `envconfig:"PASSWORD"`
	User       string        `envconfig:"USER"`
	DB         int           `envconfig:"DB"`
	MaxRetries int           `envconfig:"MAX_RETRIES"`
	CacheTTL   time.Duration `envconfig:"CACHE_TTL"`
}

func NewConfig() (*Config, error) {
	var config Config
	if err := envconfig.Process("REDIS", &config); err != nil {
		return nil, fmt.Errorf(
			"failed to make redis config: %w",
			err,
		)
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
