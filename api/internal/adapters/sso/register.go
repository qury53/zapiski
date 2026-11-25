package adapters_sso

import (
	"context"
	ssov1 "shared/pkg/proto/sso/v1"

	"github.com/qury53/GC_learn/api/internal/core/ports"
)

func (c *GRPCClient) Register(ctx context.Context, req ports.RegisterRequest) (ports.RegisterResponse, error) {
	resp, err := c.client.Register(ctx, &ssov1.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	})
	if err != nil {
		return ports.RegisterResponse{}, err
	}
	status, err := authStatusFromPB(resp.GetStatus())
	if err != nil {
		return ports.RegisterResponse{}, err
	}
	return ports.RegisterResponse{UserID: resp.GetUserId(), Status: status}, nil
}
