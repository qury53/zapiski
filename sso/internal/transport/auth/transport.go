package auth_transport

import (
	"context"
	"net/netip"
	"time"

	"shared/pkg/logger"
	ssov1 "shared/pkg/proto/sso/v1"

	"github.com/qury53/GC_learn/sso/internal/core/domain"
	oauth_transport "github.com/qury53/GC_learn/sso/internal/transport/oauth"
)

// AuthTransport — единственная точка регистрации sso.v1.AuthService на gRPC сервере.
// OAuth-методы живут в oauth_transport и подключаются явными форвардерами:
// два обработчика одного сервиса зарегистрировать нельзя, а встраивание
// дало бы неоднозначный селектор с UnimplementedAuthServiceServer.
type AuthTransport struct {
	log *logger.Logger
	ssov1.UnimplementedAuthServiceServer
	oauth *oauth_transport.OAuthTransport
	service AuthService
}

// GetOAuthConnectUrl проксирует вызов в OAuth транспорт.
func (t *AuthTransport) GetOAuthConnectUrl(
	ctx context.Context,
	req *ssov1.OAuthConnectRequest,
) (*ssov1.OAuthConnectResponse, error) {
	return t.oauth.GetOAuthConnectUrl(ctx, req)
}

// HandleOAuthCallback проксирует вызов в OAuth транспорт.
func (t *AuthTransport) HandleOAuthCallback(
	ctx context.Context,
	req *ssov1.OAuthCallbackRequest,
) (*ssov1.LoginResponse, error) {
	return t.oauth.HandleOAuthCallback(ctx, req)
}

type AuthService interface {
	Register(ctx context.Context, user domain.User) (int64, domain.VerificationStatus, error)
	Login(
		ctx context.Context,
		username string,
		password string,
		userAgent *string,
		ipAddr *netip.Addr,
	) (string, string, int64, error)
	Refresh(ctx context.Context, refreshToken string) (string, string, error)
	Logout(ctx context.Context, refreshToken string) error
	GetRefreshTokenTTL() time.Duration
	VerifyEmail(ctx context.Context, email string, code string) (domain.VerificationStatus, error)
}

type SMTPEmailSender interface {
	Send(ctx context.Context)
}

// NewAuthTransport собирает транспорт авторизации вместе с OAuth-методами.
// oauthService должен быть тем же сервисом, что и service: OAuth-методы
// обслуживаются из отдельного пакета, но тем же AuthService.
func NewAuthTransport(
	log *logger.Logger,
	service AuthService,
	oauthService oauth_transport.OAuthService,
) *AuthTransport {
	return &AuthTransport{
		log:     log,
		oauth:   oauth_transport.NewOAuthTransport(log, oauthService),
		service: service,
	}
}
