package auth_service

import (
	"context"
	"shared/pkg/logger"
	"time"

	"github.com/qury53/GC_learn/sso/internal/core/domain"
)

type AuthService struct {
	log               *logger.Logger
	userRepository    UserRepository
	sessionRepository SessionRepository
	encryptor         Encryptor
	smtp              Smtp
	cacheClient       CacheClient
	JWTSecret         []byte
	AccessTokenTTL    time.Duration
	RefreshTokenTTL   time.Duration
}

type Encryptor interface {
	HashToken(token string) string
	GeneratePasswordHash(val string) (string, error)
	CompareHashAndPassword(hash string, val string) error
	GenerateVerificationCode(length int) (string, error)
}

type SessionRepository interface {
	Create(ctx context.Context, s domain.Session) (domain.Session, error)
	GetByJTI(ctx context.Context, jti string) (domain.Session, error)
	GetByJTIForUpdate(ctx context.Context, jti string) (domain.Session, error)
	Rotate(ctx context.Context, jti string, params domain.RotateParams) (domain.Session, error)
	Revoke(ctx context.Context, jti string, reason string) error
	RevokeAllForUser(ctx context.Context, userID int64, reason string) error
	MarkCompromised(ctx context.Context, jti string, reason string) error
	ListActiveByUserID(ctx context.Context, userID int64) ([]domain.Session, error)
}

type UserRepository interface {
	SaveUser(ctx context.Context, user domain.User) (domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (domain.User, error)
}

type CacheClient interface {
	Get(ctx context.Context, key string, dst any) (bool, error)
	Set(ctx context.Context, key string, value any) error
	Del(ctx context.Context, key string) error
	GenerateCacheKey(args ...string) string
}

type Smtp interface {
	Send(
		ctx context.Context,
		toEmail string,
		code string,
	) error
}

func NewAuthService(
	log *logger.Logger,
	config *Config,
	userRepository UserRepository,
	smtp Smtp,
	cacheClient CacheClient,
) *AuthService {
	return &AuthService{
		log:             log,
		userRepository:  userRepository,
		smtp:            smtp,
		cacheClient:     cacheClient,
		JWTSecret:       []byte(config.JWTSecret),
		AccessTokenTTL:  config.AccessTokenTTL,
		RefreshTokenTTL: config.RefreshTokenTTL,
	}
}
