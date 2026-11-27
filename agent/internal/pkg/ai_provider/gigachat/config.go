package ai_provider_gc

import "github.com/kelseyhightower/envconfig"

type Config struct {
	APIKey string `envconfig:"API_KEY" required:"true"`
	Scope  string `envconfig:"SCOPE" default:"GIGACHAT_API_PERS"`
}

func NewConfig() (*Config, error) {
	var config Config
	if err := envconfig.Process("GIGACHAT", &config); err != nil {
		return nil, err
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
