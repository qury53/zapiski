package auth_service

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config конфиг сервиса авторизации

type Config struct {
	JWTSecret       string        `envconfig:"JWT_SECRET" required:"true"`
	RefreshTokenTTL time.Duration `envconfig:"REFRESH_TOKEN_TTL" default:"720h"`
	AccessTokenTTL  time.Duration `envconfig:"ACCESS_TOKEN_TTL" default:"900s"`

	GoogleClientID           string        `envconfig:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret       string        `envconfig:"GOOGLE_CLIENT_SECRET"`
	OAuthRedirectURI         string        `envconfig:"OAUTH_REDIRECT_URI"`
	OAuthStateTTL            time.Duration `envconfig:"OAUTH_STATE_TTL" default:"10m"`
	OAuthAllowedRedirectURIs []string      `envconfig:"OAUTH_ALLOWED_REDIRECT_URIS"`
}

// NewConfig создат конфиг для AuthService
//
// Returns: *Config, error

func NewConfig() (*Config, error) {
	var config Config
	if err := envconfig.Process("AUTH", &config); err != nil {
		return nil, fmt.Errorf("failed to proccess auth config: %w", err)
	}

	return &config, nil
}

// NewConfigMust создает конфиг AuthService, использует NewConfig,
// при ошибке вызывает панику
//
// Returns: *Config

func NewConfigMust() *Config {
	config, err := NewConfig()
	if err != nil {
		panic(err)
	}

	return config
}
