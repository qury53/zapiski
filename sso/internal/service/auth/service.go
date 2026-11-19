package auth_service

import (
	"context"
	"time"

	"github.com/qury53/GC_learn/platform/logger"
	"github.com/qury53/GC_learn/sso/internal/core/domain"
	core_oauth "github.com/qury53/GC_learn/sso/internal/core/oauth"
	"github.com/redis/go-redis/v9"
)

type AuthService struct {
	log               *logger.Logger
	userRepository    UserRepository
	sessionRepository SessionRepository
	oauthRepository   OAuthRepository
	oauthClients      map[domain.Provider]core_oauth.Client
	allowedRedirects  map[string]struct{}
	oauthStateTTL     time.Duration
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
	GenerateVerificationCode() (string, error)
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

type OAuthRepository interface {
	GetByProviderUserID(ctx context.Context, provider domain.Provider, providerUserID string) (domain.OAuthAccount, error)
	Create(ctx context.Context, account domain.OAuthAccount) (domain.OAuthAccount, error)
	Touch(ctx context.Context, id int64, email string, emailVerified bool) error
}

type UserRepository interface {
	SaveUser(ctx context.Context, user domain.User) (domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	GetUserByID(ctx context.Context, id int64) (domain.User, error)
	CreateOAuthUser(ctx context.Context, username string, email string) (domain.User, error)
	ActivateUser(ctx context.Context, email string) error
}

type CacheClient interface {
	redis.Scripter
	Get(ctx context.Context, key string, dst any) (bool, error)
	Set(ctx context.Context, key string, value any) error
	SetWithTTL(ctx context.Context, key string, value any, ttl time.Duration) error
	Incr(ctx context.Context, key string) (int64, error)
	Del(ctx context.Context, key string) error
	GetDel(ctx context.Context, key string, dst any) (bool, error)
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
	sessionRepository SessionRepository,
	oauthRepository OAuthRepository,
	encryptor Encryptor,
	smtp Smtp,
	cacheClient CacheClient,
	oauthClients map[domain.Provider]core_oauth.Client,
) *AuthService {
	allowedRedirects := make(map[string]struct{}, len(config.OAuthAllowedRedirectURIs))
	for _, uri := range config.OAuthAllowedRedirectURIs {
		allowedRedirects[uri] = struct{}{}
	}

	return &AuthService{
		log:               log,
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
		oauthRepository:   oauthRepository,
		oauthClients:      oauthClients,
		allowedRedirects:  allowedRedirects,
		oauthStateTTL:     config.OAuthStateTTL,
		encryptor:         encryptor,
		smtp:              smtp,
		cacheClient:       cacheClient,
		JWTSecret:         []byte(config.JWTSecret),
		AccessTokenTTL:    config.AccessTokenTTL,
		RefreshTokenTTL:   config.RefreshTokenTTL,
	}
}
