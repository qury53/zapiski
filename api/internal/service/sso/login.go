package sso

import (
	"context"
	"fmt"

	"github.com/qury53/GC_learn/api/internal/core/ports"
	transport_sso "github.com/qury53/GC_learn/api/internal/transport/sso"
)

func (s *Service) Login(ctx context.Context, req transport_sso.LoginRequest, headers ports.HeaderCarrier) (transport_sso.LoginResponse, error) {
	resp, err := s.client.Login(ctx, ports.LoginRequest{Username: req.Username, Password: req.Password}, headers)
	if err != nil {
		return transport_sso.LoginResponse{}, fmt.Errorf("login: %w", err)
	}
	return transport_sso.LoginResponse{AccessToken: resp.AccessToken, UserID: resp.UserID}, nil
}
