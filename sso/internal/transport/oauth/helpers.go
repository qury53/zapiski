package oauth_transport

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/netip"
	"strings"
	"time"

	ssov1 "shared/pkg/proto/sso/v1"

	"github.com/qury53/GC_learn/sso/internal/core/domain"
	core_errors "github.com/qury53/GC_learn/sso/internal/core/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// mapProvider переводит provider из proto в доменный тип.
// Неизвестный provider — ошибка клиента, поэтому оборачиваем в ErrInvalidArgument,
// иначе RespondCode отдал бы Internal.
//
// Returns: domain.Provider, error
func mapProvider(provider ssov1.OAuthProvider) (domain.Provider, error) {
	var result domain.Provider

	switch provider {
	case ssov1.OAuthProvider_OAUTH_PROVIDER_GOOGLE:
		result = domain.ProviderGoogle
	case ssov1.OAuthProvider_OAUTH_PROVIDER_GITHUB:
		result = domain.ProviderGithub
	case ssov1.OAuthProvider_OAUTH_PROVIDER_YANDEX:
		result = domain.ProviderYandex
	default:
		return "", fmt.Errorf("%w: unknown provider %s", core_errors.ErrInvalidArgument, provider)
	}

	if _, ok := domain.SupportedProviders[result]; !ok {
		return "", fmt.Errorf("%w: %s", core_errors.ErrInvalidArgument, result)
	}

	return result, nil
}

// setRefreshCookie кладет refresh токен в HttpOnly cookie, как это делает Login.
//
// Returns: error
func (t *OAuthTransport) setRefreshCookie(
	ctx context.Context,
	refreshToken string,
) error {
	cookieValue := fmt.Sprintf(
		"refresh_token=%s; Path=/; Max-Age=%d; HttpOnly; Secure; SameSite=Lax",
		refreshToken,
		int(t.service.GetRefreshTokenTTL()/time.Second),
	)

	if err := grpc.SendHeader(ctx, metadata.Pairs("set-cookie", cookieValue)); err != nil {
		return fmt.Errorf("failed to set cookie: %w", err)
	}

	return nil
}

// getUserAgent достает user-agent из метаданных запроса.
//
// Returns: string
func getUserAgent(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	if userAgent := md.Get("user-agent"); len(userAgent) > 0 {
		return userAgent[0]
	}

	return ""
}

// getClientIP достает ip клиента из метаданных прокси или из peer.
//
// Returns: *netip.Addr
func getClientIP(ctx context.Context) *netip.Addr {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if xff := md.Get("x-forwarded-for"); len(xff) > 0 {
			rawIP, _, _ := strings.Cut(xff[0], ",")

			if addr, err := netip.ParseAddr(strings.TrimSpace(rawIP)); err == nil {
				unmapped := addr.Unmap()

				return &unmapped
			}
		}

		if xri := md.Get("x-real-ip"); len(xri) > 0 {
			if addr, err := netip.ParseAddr(strings.TrimSpace(xri[0])); err == nil {
				unmapped := addr.Unmap()

				return &unmapped
			}
		}
	}

	p, ok := peer.FromContext(ctx)
	if !ok || p.Addr == nil {
		return nil
	}

	if tcpAddr, ok := p.Addr.(*net.TCPAddr); ok {
		if addr, ok := netip.AddrFromSlice(tcpAddr.IP); ok {
			unmapped := addr.Unmap()

			return &unmapped
		}
	}

	if addrPort, err := netip.ParseAddrPort(p.Addr.String()); err == nil {
		unmapped := addrPort.Addr().Unmap()

		return &unmapped
	}

	return nil
}

// logOAuthError логирует ошибку OAuth. code, state и code_verifier сюда не попадают.
func (t *OAuthTransport) logOAuthError(
	msg string,
	provider domain.Provider,
	err error,
) {
	t.log.Debug(msg,
		slog.String("provider", string(provider)),
		slog.Any("error", err),
	)
}
