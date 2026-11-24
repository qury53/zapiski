package email_smtp_repository

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Username                string `envconfig:"USERNAME"`
	Password                string `envconfig:"PASSWORD"`
	SenderEmail             string `envconfig:"SENDER_EMAIL" required:"true"`
	Addr                    string `envconfig:"ADDR" required:"true"`
	Host                    string `envconfig:"HOST" required:"true"`
	Port                    string `envconfig:"PORT" required:"true"`
	VerificationCodeTimeout string `envconfig:"VERIFICATION_CODE_TIMEOUT" default:"300s"`
}

func NewConfig() (*Config, error) {
	var config Config
	if err := envconfig.Process("SMTP", &config); err != nil {
		return nil, fmt.Errorf("failed to process .env: %w", err)
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
