package sso

import (
	"context"
	"fmt"

	"github.com/qury53/GC_learn/api/internal/core/ports"
	transport_sso "github.com/qury53/GC_learn/api/internal/transport/sso"
)

func (s *Service) Refresh(ctx context.Context, refreshToken string, headers ports.HeaderCarrier) (transport_sso.TokenPair, error) {
	resp, err := s.client.Refresh(ctx, refreshToken, headers)
	if err != nil {
		return transport_sso.TokenPair{}, fmt.Errorf("refresh: %w", err)
	}
	return transport_sso.TokenPair{AccessToken: resp.AccessToken, RefreshToken: resp.RefreshToken}, nil
}
