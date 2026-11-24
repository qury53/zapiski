package adapters_sso

import (
	"context"
	ssov1 "shared/pkg/proto/sso/v1"

	"github.com/qury53/GC_learn/api/internal/core/ports"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (c *GRPCClient) VerifyEmail(ctx context.Context, req ports.VerifyEmailRequest, headers ports.HeaderCarrier) (ports.VerifyEmailResponse, error) {
	var headerMD, trailerMD metadata.MD
	resp, err := c.client.VerifyEmail(ctx, &ssov1.VerifyEmailRequest{Email: req.Email, Code: req.Code}, grpc.Header(&headerMD), grpc.Trailer(&trailerMD))
	copySetCookies(headers, headerMD, trailerMD)
	if err != nil {
		return ports.VerifyEmailResponse{}, err
	}
	status, err := authStatusFromPB(resp.GetStatus())
	if err != nil {
		return ports.VerifyEmailResponse{}, err
	}
	return ports.VerifyEmailResponse{Status: status}, nil
}
