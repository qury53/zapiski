package auth

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// Config конфиг валидации JWT.
// Читает те же переменные, что и sso (AUTH_JWT_SECRET / AUTH_JWT_ISSUER).
type Config struct {
	Secret string `envconfig:"JWT_SECRET" required:"true"`
	Issuer string `envconfig:"JWT_ISSUER" default:"sso-zapiski"`
}

func NewConfig() (*Config, error) {
	var config Config
	if err := envconfig.Process("AUTH", &config); err != nil {
		return nil, fmt.Errorf("failed to process auth config: %w", err)
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
