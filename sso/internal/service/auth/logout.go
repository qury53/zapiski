package auth_service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/qury53/GC_learn/sso/internal/core/domain"
)

func (s *AuthService) Logout(
	ctx context.Context,
	refreshToken string,
) error {
	claims, err := s.ParseRefreshToken(refreshToken)
	if err != nil {
		return fmt.Errorf("auth.Logout: failed to parse refresh token: %w", err)
	}

	if err := s.sessionRepository.Revoke(
		ctx,
		claims.ID,
		domain.RevokedReasonLogout,
	); err != nil {
		s.log.Info(
			"auth.Logout revoke err",
			slog.Int64("userID", claims.UserID),
			slog.Any("error", err),
		)
	}

	return nil
}
