package oauth_google_repository

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ClientID     string `envconfig:"CLIENT_ID" required:"true"`
	ClientSecret string `envconfig:"CLIENT_SECRET" required:"true"`
	RedirectURI  string `envconfig:"REDIRECT_URI" required:"true"`
}

func NewConfig() (*Config, error) {
	var config Config
	if err := envconfig.Process("OAUTH_GOOGLE", &config); err != nil {
		return nil, fmt.Errorf("process oauth google config: %w", err)
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
