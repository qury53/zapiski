package auth_service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/qury53/GC_learn/sso/internal/core/domain"
	core_errors "github.com/qury53/GC_learn/sso/internal/core/errors"
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

		// Logout идемпотентен: неизвестная или уже отозванная сессия —
		// не ошибка для клиента, который просто хочет выйти.
		if errors.Is(err, core_errors.ErrNotFound) || errors.Is(err, core_errors.ErrSessionRevoked) {
			return nil
		}

		return fmt.Errorf("auth.Logout: revoke session: %w", err)
	}

	return nil
}
