package auth_transport

import (
	"context"
	"net/netip"
	"shared/pkg/logger"
	ssov1 "shared/pkg/proto/sso/v1"
	"time"

	"github.com/qury53/GC_learn/sso/internal/core/domain"
)

type AuthTransport struct {
	log *logger.Logger
	ssov1.UnimplementedAuthServiceServer
	service AuthService
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
	VerifyEmail(ctx context.Context, email string, code string) domain.VerificationStatus
	GetConnectURL(
		ctx context.Context,
		provider domain.Provider,
		redirectURI string,
		challenge string,
	) (url string, state string, err error)
	ExchangeCallback(
		ctx context.Context,
		provider domain.Provider,
		code string,
		state string,
		codeVerifier string,
	)
}

type SMTPEmailSender interface {
	Send(ctx context.Context)
}
