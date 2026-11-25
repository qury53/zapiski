package logger

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// Config структура конфига логгера
type Config struct {
	LogLevel int `envconfig:"LOG_LEVEL" default:"-4"`
}

// NewConfig парсит .env файл и создает конфиг логгера
//
// Returns: Config, error
func NewConfig() (*Config, error) {
	var config Config
	if err := envconfig.Process("LOGGER", &config); err != nil {
		return nil, fmt.Errorf("failed to proccess logger config: %w", err)
	}

	return &config, nil
}

// NewConfigMust создает конфиг логгера, при ошибки вызывает панику
//
// Returns: Config

func NewConfigMust() *Config {
	config, err := NewConfig()
	if err != nil {
		panic(err)
	}

	return config
}
