package adapters_sso

import (
	"context"
	ssov1 "shared/pkg/proto/sso/v1"

	"github.com/qury53/GC_learn/api/internal/core/ports"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (c *GRPCClient) Login(ctx context.Context, req ports.LoginRequest, headers ports.HeaderCarrier) (ports.LoginResponse, error) {
	var headerMD, trailerMD metadata.MD
	resp, err := c.client.Login(ctx, &ssov1.LoginRequest{Username: req.Username, Password: req.Password}, grpc.Header(&headerMD), grpc.Trailer(&trailerMD))
	copySetCookies(headers, headerMD, trailerMD)
	if err != nil {
		return ports.LoginResponse{}, err
	}
	return ports.LoginResponse{AccessToken: resp.GetAccessToken(), UserID: resp.GetUserId()}, nil
}
