package user_service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/qury53/GC_learn/user/internal/core/domain"
	core_errors "github.com/qury53/GC_learn/user/internal/core/errors"
)

const (
	maxUsers = 200
)

func (s *UserService) GetUsers(
	ctx context.Context,
	limit *int64,
	offset *int64,
) ([]domain.UserProfile, error) {
	if limit != nil && (*limit <= 0 || *limit > maxUsers) {
		return nil, fmt.Errorf(
			"GetUsers service: limit has to be positive if set: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf(
			"GetUsers service: offset has to be non-negative if set: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	var users []domain.UserProfile
	cacheKey := s.cachClient.GenerateCacheKey("profile")
	if limit != nil {
		cacheKey = s.cachClient.GenerateCacheKey(cacheKey, fmt.Sprintf("%d", *limit))
	}

	if offset != nil {
		cacheKey = s.cachClient.GenerateCacheKey(cacheKey, fmt.Sprintf("%d", *offset))
	}

	hit, err := s.cachClient.Get(ctx, cacheKey, &users)
	if err != nil {
		s.log.Warn(
			"cache error",
			slog.String("method", "GetUsers"),
			slog.String("cacheKey", cacheKey),
			slog.Any("error", err),
		)

		hit = false
	}

	if hit {
		return users, nil
	}

	users, err = s.userRepository.GetUsers(
		ctx,
		limit,
		offset,
	)

	if err != nil {
		return nil, fmt.Errorf("GetUsers service: %w", err)
	}

	if err := s.cachClient.Set(ctx, cacheKey, users); err != nil {
		s.log.Warn(
			"cache error",
			slog.String("method", "GetUsers"),
			slog.String("cacheKey", cacheKey),
			slog.Any("error", err),
		)
	}

	return users, nil
}
