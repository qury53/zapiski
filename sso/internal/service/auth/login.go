package auth_service

import (
	"context"
	"fmt"
	"net/netip"
	"time"

	"github.com/google/uuid"
	"github.com/qury53/GC_learn/sso/internal/core/domain"
	core_errors "github.com/qury53/GC_learn/sso/internal/core/errors"
)

func (s *AuthService) Login(
	ctx context.Context,
	username string,
	password string,
	userAgent *string,
	ipAddr *netip.Addr,
) (string, string, int64, error) {
	user, err := s.userRepository.GetUserByUsername(ctx, username)
	if err != nil {
		return "", "", 0, fmt.Errorf(
			"service: get user by username: %w",
			err,
		)
	}

	if !user.CanAuthorize() {
		s.log.Debug("login rejected: user not authorized", "username", username, "state", user.State)
		return "", "", 0, core_errors.ErrInvalidCredentials
	}

	if err := s.encryptor.CompareHashAndPassword(
		user.Password,
		password,
	); err != nil {
		return "", "", 0, fmt.Errorf(
			"service: invalid credentials: %w",
			err,
		)
	}

	jti := uuid.NewString()
	now := time.Now()
	expiresAt := now.Add(s.RefreshTokenTTL)

	refreshToken, err := s.GenerateRefreshToken(user.ID, jti)
	if err != nil {
		return "", "", 0, fmt.Errorf("service: generate refresh token: %w", err)
	}

	session := domain.NewSessionUninitialized(
		user.ID,
		jti,
		s.encryptor.HashToken(refreshToken),
		userAgent,
		ipAddr,
		expiresAt,
		nil,
	)

	if _, err = s.sessionRepository.Create(ctx, session); err != nil {
		return "", "", 0, fmt.Errorf(
			"failed create session: %w",
			err,
		)
	}

	accessToken, err := s.GenerateAccessToken(user.ID)
	if err != nil {
		return "", "", 0, fmt.Errorf("service: generate access token: %w", err)
	}

	return accessToken, refreshToken, user.ID, nil
}
