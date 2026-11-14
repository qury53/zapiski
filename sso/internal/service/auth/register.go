package auth_service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/qury53/GC_learn/sso/internal/core/domain"
	core_encrypt "github.com/qury53/GC_learn/sso/internal/core/encrypt"
)

func (s *AuthService) Register(
	ctx context.Context,
	user domain.User,
) (int64, domain.VerificationStatus, error) {
	hash, err := core_encrypt.GeneratePasswordHash(user.Password)
	if err != nil {
		s.log.Debug("failed to generate password hash", slog.String("error", err.Error()))
		return 0, domain.StatusUnknown, fmt.Errorf(
			"generate password hash: %w",
			err,
		)
	}

	user.Password = hash
	user, err = s.userRepository.SaveUser(ctx, user)
	if err != nil {
		return 0, domain.StatusUnknown, fmt.Errorf(
			"save user to db: %w",
			err,
		)
	}

	s.log.Debug("creating verification code")
	verificationCode, err := s.GenerateVerificationCode()

	if err != nil {
		return 0, domain.StatusUnknown, fmt.Errorf(
			"generate verification code: %w",
			err,
		)
	}

	if err := s.CacheVerificationCode(
		ctx,
		user.Email,
		verificationCode,
	); err != nil {
		s.log.Warn("failed to cache verification code", slog.String("error", err.Error()))
		return 0, domain.StatusUnknown, fmt.Errorf("failed to cache verification code: %w", err)
	}

	if err := s.smtp.Send(
		ctx,
		user.Email,
		verificationCode,
	); err != nil {
		return 0, domain.StatusUnknown, fmt.Errorf(
			"send verification code: %w",
			err,
		)
	}

	return user.ID, domain.StatusPendingVerification, nil
}
