package sso

import (
	"context"
	"fmt"

	"github.com/qury53/GC_learn/api/internal/core/ports"
	transport_sso "github.com/qury53/GC_learn/api/internal/transport/sso"
)

func (s *Service) VerifyEmail(ctx context.Context, req transport_sso.VerifyEmailRequest, headers ports.HeaderCarrier) (transport_sso.VerifyEmailResponse, error) {
	resp, err := s.client.VerifyEmail(ctx, ports.VerifyEmailRequest{Email: req.Email, Code: req.Code}, headers)
	if err != nil {
		return transport_sso.VerifyEmailResponse{}, fmt.Errorf("verify email: %w", err)
	}
	return transport_sso.VerifyEmailResponse{Status: transport_sso.AuthStatus(resp.Status)}, nil
}
