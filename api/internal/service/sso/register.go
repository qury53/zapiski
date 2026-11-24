package sso

import (
	"context"
	"fmt"

	"github.com/qury53/GC_learn/api/internal/core/ports"
	transport_sso "github.com/qury53/GC_learn/api/internal/transport/sso"
)

func (s *Service) Register(ctx context.Context, req transport_sso.RegisterRequest) (transport_sso.RegisterResponse, error) {
	resp, err := s.client.Register(ctx, ports.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	})
	if err != nil {
		return transport_sso.RegisterResponse{}, fmt.Errorf("register: %w", err)
	}

	return transport_sso.RegisterResponse{UserID: resp.UserID, Status: string(resp.Status)}, nil
}
