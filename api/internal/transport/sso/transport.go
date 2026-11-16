package sso

import (
	"shared/pkg/logger"
	ssov1 "shared/pkg/proto/sso/v1"
)

type SSOTransport struct {
	log   *logger.Logger
	ssov1 ssov1.AuthServiceClient
}

func NewSSOTransport(log *logger.Logger, sso ssov1.AuthServiceClient) *SSOTransport {
	return &SSOTransport{
		log:   log,
		ssov1: sso,
	}
}
