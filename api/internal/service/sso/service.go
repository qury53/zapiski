package sso

import (
	"github.com/qury53/GC_learn/api/internal/core/ports"
	"github.com/qury53/GC_learn/platform/logger"
)

type Service struct {
	client ports.SSOClient
	log    *logger.Logger
}

func NewAuthService(client ports.SSOClient, log *logger.Logger) *Service {
	return &Service{
		client: client,
		log:    log,
	}
}
