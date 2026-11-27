package user_service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/qury53/GC_learn/user/internal/core/domain"
	core_errors "github.com/qury53/GC_learn/user/internal/core/errors"
)

func (s *UserService) GetUserByID(
	ctx context.Context,
	userID int64,
) (domain.UserProfile, error) {
	if userID < 1 {
		return domain.UserProfile{}, fmt.Errorf(
			"GetUserByID service: userID has to be positive: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	cacheKey := s.cachClient.GenerateCacheKey("profile", fmt.Sprintf("%d", userID))
	//ищем в кэше
	var user domain.UserProfile
	hit, err := s.cachClient.Get(
		ctx,
		cacheKey,
		&user,
	)

	if err != nil {
		s.log.Warn(
			"cache error",
			slog.String("method", "GetUserByID"),
			slog.String("cacheKey", cacheKey),
			slog.Any("error", err),
		)

		hit = false
	}

	if hit {
		s.log.Debug(
			"cache hit",
			slog.Int64("userID", userID),
		)

		return user, nil
	}

	user, err = s.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		return domain.UserProfile{}, fmt.Errorf(
			"GetUserByID service: %w",
			err,
		)
	}

	//пишем в кэш
	if err := s.cachClient.Set(
		ctx,
		cacheKey,
		user,
	); err != nil {
		s.log.Warn(
			"failed to set cache",
			slog.String("method", "GetUserByID"),
			slog.Any("error", err),
		)
	}

	return user, nil
}
