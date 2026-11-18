package auth_service

import (
	"context"
	"fmt"
	"net/netip"

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

	return s.issueSession(ctx, user.ID, userAgent, ipAddr)
}
