package ai_provider_smartapi

import "net/http"

type SmartAPIProvider struct {
	client    *http.Client
	secretKey string
}

func NewSmartAPIProvider(config *Config) *SmartAPIProvider {
	return &SmartAPIProvider{
		client:    &http.Client{},
		secretKey: config.SecretKey,
	}
}
