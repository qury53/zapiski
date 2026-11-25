package adapters_sso

import (
	"context"
	ssov1 "shared/pkg/proto/sso/v1"

	"github.com/qury53/GC_learn/api/internal/core/ports"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (c *GRPCClient) Refresh(ctx context.Context, refreshToken string, headers ports.HeaderCarrier) (ports.TokenPair, error) {
	var headerMD, trailerMD metadata.MD
	resp, err := c.client.Refresh(ctx, &ssov1.RefreshRequest{RefreshToken: refreshToken}, grpc.Header(&headerMD), grpc.Trailer(&trailerMD))
	copySetCookies(headers, headerMD, trailerMD)
	if err != nil {
		return ports.TokenPair{}, err
	}
	return ports.TokenPair{AccessToken: resp.GetAccessToken(), RefreshToken: resp.GetRefreshToken()}, nil
}
