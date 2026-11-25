package auth_service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/qury53/GC_learn/sso/internal/core/domain"
)

const (
	codeLength = 6
)

func (s *AuthService) Register(
	ctx context.Context,
	user domain.User,
) (int64, domain.VerificationStatus, error) {
	// Email в БД хранится в нижнем регистре: на нём завязаны и UNIQUE(email),
	// и точное сравнение в ActivateUser/GetUserByEmail, и OAuth-автопривязка.
	user.Email = normalizeEmail(user.Email)

	hash, err := s.encryptor.GeneratePasswordHash(user.Password)
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

	// Ошибка на любом следующем шаге оставит пользователя в pending_verification
	// без кода, а повторный Register упадёт на UNIQUE(email): откатываем всё.
	cleanup := func() {
		cleanupCtx := context.WithoutCancel(ctx)

		if dropErr := s.dropVerificationCode(cleanupCtx, user.Email); dropErr != nil {
			s.log.Warn("failed to drop verification code", slog.String("error", dropErr.Error()))
		}

		if delErr := s.userRepository.DeletePendingUser(cleanupCtx, user.ID); delErr != nil {
			s.log.Warn("failed to delete pending user", slog.String("error", delErr.Error()))
		}
	}

	s.log.Debug("creating verification code")
	verificationCode, err := s.encryptor.GenerateVerificationCode()
	if err != nil {
		cleanup()
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
		cleanup()
		return 0, domain.StatusUnknown, fmt.Errorf("failed to cache verification code: %w", err)
	}

	if err := s.smtp.Send(
		ctx,
		user.Email,
		verificationCode,
	); err != nil {
		cleanup()
		return 0, domain.StatusUnknown, fmt.Errorf(
			"send verification code: %w",
			err,
		)
	}

	return user.ID, domain.StatusPendingVerification, nil
}
