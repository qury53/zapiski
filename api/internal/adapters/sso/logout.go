package adapters_sso

import (
	"context"
	ssov1 "shared/pkg/proto/sso/v1"

	"github.com/qury53/GC_learn/api/internal/core/ports"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (c *GRPCClient) Logout(ctx context.Context, refreshToken string, headers ports.HeaderCarrier) error {
	var headerMD, trailerMD metadata.MD
	_, err := c.client.Logout(ctx, &ssov1.LogoutRequest{RefreshToken: refreshToken}, grpc.Header(&headerMD), grpc.Trailer(&trailerMD))
	copySetCookies(headers, headerMD, trailerMD)
	return err
}
