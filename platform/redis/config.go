package redis

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Addr       string        `envconfig:"ADDR" required:"true"`
	Password   string        `envconfig:"PASSWORD" required:"true"`
	User       string        `envconfig:"USER" required:"true"`
	DB         int           `envconfig:"DB" required:"true"`
	MaxRetries int           `envconfig:"MAX_RETRIES" required:"true"`
	CacheTTL   time.Duration `envconfig:"CACHE_TTL" required:"true"`
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
