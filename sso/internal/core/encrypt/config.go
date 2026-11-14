package core_encrypt

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// Config конфиг для encrypt-ора
type Config struct {
	EncryptionSecret string `envconfig:"ENCRYPTION_SECRET" required:"true"`
}

// NewConfig создает конфиг для encrypt-ора
//
// Returns: *Config, error

func NewConfig() (*Config, error) {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		return nil, fmt.Errorf("encryption config process error: %w", err)
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
