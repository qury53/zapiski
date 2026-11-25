package adapters_sso

import (
	ssov1 "shared/pkg/proto/sso/v1"

	"github.com/qury53/GC_learn/platform/logger"
)

type GRPCClient struct {
	client ssov1.AuthServiceClient
	log    *logger.Logger
}

func NewGRPCClient(client ssov1.AuthServiceClient, log *logger.Logger) *GRPCClient {
	return &GRPCClient{client: client, log: log}
}
