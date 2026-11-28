package user_transport

import (
	"context"
	userv1 "shared/pkg/proto/user/v1"

	"github.com/qury53/GC_learn/platform/logger"
	"github.com/qury53/GC_learn/user/internal/core/domain"
)

type UserTransport struct {
	userv1.UnimplementedUserServer
	service UserService
	log     *logger.Logger
}

type UserService interface {
	GetUsers(
		ctx context.Context,
		limit *int64,
		offset *int64,
	) ([]domain.UserProfile, error)

	GetUserByID(
		ctx context.Context,
		userID int64,
	) (domain.UserProfile, error)
}

func NewUserTransport(log *logger.Logger, service UserService) *UserTransport {
	return &UserTransport{
		service: service,
		log:     log,
	}
}
