package user_service

import (
	"context"

	"github.com/qury53/GC_learn/platform/logger"
	"github.com/qury53/GC_learn/user/internal/core/domain"
)

type UserService struct {
	log            *logger.Logger
	userRepository UserRepository
	cachClient     CacheClient
}

type UserRepository interface {
	GetUsers(
		ctx context.Context,
		limit *int64,
		offset *int64,
	) ([]domain.UserProfile, error)

	GetUserByID(
		ctx context.Context,
		user_id int64,
	) (domain.UserProfile, error)
}

type CacheClient interface {
	Get(ctx context.Context, key string, dst any) (bool, error)
	Set(ctx context.Context, key string, value any) error
	Del(ctx context.Context, key string) error
	GenerateCacheKey(args ...string) string
}

func NewUserService(
	log *logger.Logger,
	userRepository UserRepository,
	cacheClient CacheClient,
) *UserService {
	return &UserService{
		log:            log,
		userRepository: userRepository,
		cachClient:     cacheClient,
	}
}
