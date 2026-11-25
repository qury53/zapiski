package ai_provider_smartapi

import (
	"net/http"
	"time"

	"github.com/qury53/GC_learn/platform/logger"
)

type SmartAPIProvider struct {
	client    *http.Client
	log       *logger.Logger
	secretKey string
}

func NewSmartAPIProvider(config *Config, log *logger.Logger) *SmartAPIProvider {
	return &SmartAPIProvider{
		client:    &http.Client{Timeout: 90 * time.Second},
		log:       log,
		secretKey: config.SecretKey,
	}
}
