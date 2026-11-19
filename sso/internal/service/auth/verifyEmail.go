package auth_service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/qury53/GC_learn/sso/internal/core/domain"
	core_errors "github.com/qury53/GC_learn/sso/internal/core/errors"
)

const (
	maxAttempts = 3
)
const (
	resOK       int64 = 1
	resWrong    int64 = 0
	resNotFound int64 = -1
	resLocked   int64 = -2
)

func (s *AuthService) VerifyEmail(
	ctx context.Context,
	email string,
	code string,
) (domain.VerificationStatus, error) {
	res, err := verifyEmailScript.Run(ctx, s.cacheClient,
		[]string{s.verificationCacheKey(email), s.verificationCodeAttemptsCacheKey(email)},
		code,
		maxAttempts,
		verificationCodeTTL.Milliseconds(),
	).Int64()
	if err != nil {
		s.log.ErrorContext(ctx, "verify: script failed", slog.String("error", err.Error()))
		return domain.StatusUnknown, fmt.Errorf("verify script: %w", err) // 500, не 401
	}

	switch res {
	case resOK:
		if err := s.userRepository.ActivateUser(ctx, email); err != nil {
			s.log.ErrorContext(ctx, "verify: activate failed", slog.String("error", err.Error()))
			return domain.StatusUnknown, fmt.Errorf("activate user: %w", err)
		}
		return domain.StatusActive, nil

	case resWrong:
		return domain.StatusUnknown, core_errors.ErrUnauthorized

	case resLocked:
		return domain.StatusUnknown, core_errors.ErrTooManyRequest

	default: // resNotFound
		return domain.StatusUnknown, core_errors.ErrUnauthorized
	}
}
