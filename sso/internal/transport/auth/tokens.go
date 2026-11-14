package auth_transport

import (
	"context"
	"log/slog"
	"net/http"
	"shared/pkg/logger"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// extractRefreshFromCookie достает refresh токен из cookie
// если куки нет, то возвращается пустая строка
// Returns: string

func extractRefreshFromCookie(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	var rawCookies []string
	rawCookies = append(rawCookies, md.Get("grpcgateway-cookie")...)
	rawCookies = append(rawCookies, md.Get("cookie")...)

	if len(rawCookies) == 0 {
		return ""
	}

	req := &http.Request{Header: http.Header{}}
	for _, line := range rawCookies {
		req.Header.Add("Cookie", strings.TrimSpace(line))
	}

	if cookie, err := req.Cookie("refresh_token"); err == nil {
		return cookie.Value
	}
	return ""
}

// clearRefreshCookie очищает refresh_token куки если истек срок жизни
//
// Returns: error

func clearRefreshCookie(ctx context.Context, log *logger.Logger) {
	cookie := "refresh_token" + "=; Path=/; Max-Age=0; HttpOnly; Secure; SameSite=Lax"
	if err := grpc.SendHeader(ctx, metadata.Pairs("set-cookie", cookie)); err != nil {
		log.Warn("failed to clear refresh cookie", slog.Any("error", err))
	}
}
