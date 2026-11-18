package pgx_pool

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Конфиг Pool
type Config struct {
	Host     string        `envconfig:"HOST" required:"true"`
	Port     string        `envconfig:"PORT" required:"true"`
	DB       string        `envconfig:"DB" required:"true"`
	User     string        `envconfig:"USER" required:"true"`
	Password string        `envconfig:"PASSWORD" required:"true"`
	Timeout  time.Duration `envconfig:"TIMEOUT" default:"15s"`
}

// Создает конфиг для пула подключения
func NewConfig() (*Config, error) {
	var config Config
	if err := envconfig.Process("POSTGRES", &config); err != nil {
		return nil, fmt.Errorf("failed to process postgres .env: %w", err)
	}

	return &config, nil
}

// Создает конфиг для пула подключения, при неудаче вызывает панику
func NewConfigMust() *Config {
	config, err := NewConfig()
	if err != nil {
		panic(err)
	}

	return config
}
