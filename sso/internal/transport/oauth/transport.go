package oauth_transport

import (
	"context"
	"net/netip"
	"time"

	"github.com/qury53/GC_learn/platform/logger"
	"github.com/qury53/GC_learn/sso/internal/core/domain"
)

// OAuthTransport реализует OAuth-методы sso.v1.AuthService:
// GetOAuthConnectUrl и HandleOAuthCallback.
// Не регистрируется отдельно от AuthTransport: иначе на одном gRPC сервере
// оказалось бы два обработчика одного и того же сервиса.
type OAuthTransport struct {
	log     *logger.Logger
	service OAuthService
}

type OAuthService interface {
	// GetConnectURI собирает URL согласия провайдера и возвращает state,
	// который сервис сохранил для проверки в callback.
	GetConnectURI(
		ctx context.Context,
		provider domain.Provider,
		redirectURI string,
		challenge string,
	) (authorizationURL string, state string, err error)

	// ExchangeCallback меняет authorization code на пару локальных токенов.
	ExchangeCallback(
		ctx context.Context,
		provider domain.Provider,
		code string,
		state string,
		codeVerifier string,
		userAgent *string,
		ipAddr *netip.Addr,
	) (accessToken string, refreshToken string, userID int64, err error)

	GetRefreshTokenTTL() time.Duration
}

func NewOAuthTransport(
	log *logger.Logger,
	service OAuthService,
) *OAuthTransport {
	return &OAuthTransport{
		log:     log,
		service: service,
	}
}
