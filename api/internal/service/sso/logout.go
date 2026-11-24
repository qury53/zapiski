package sso

import (
	"context"
	"fmt"

	"github.com/qury53/GC_learn/api/internal/core/ports"
)

func (s *Service) Logout(ctx context.Context, refreshToken string, headers ports.HeaderCarrier) error {
	if err := s.client.Logout(ctx, refreshToken, headers); err != nil {
		return fmt.Errorf("logout: %w", err)
	}
	return nil
}
