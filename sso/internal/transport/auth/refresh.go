package auth_transport

import (
	"context"
	"fmt"
	"log/slog"
	ssov1 "shared/pkg/proto/sso/v1"
	"time"

	core_transport_grpc_response "github.com/qury53/GC_learn/sso/internal/core/transport/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (t *AuthTransport) Refresh(
	ctx context.Context,
	req *ssov1.RefreshRequest,
) (*ssov1.RefreshResponse, error) {
	refreshToken := extractRefreshFromCookie(ctx)
	fromCookie := refreshToken != ""
	if !fromCookie {
		refreshToken = req.GetRefreshToken()
	}

	if refreshToken == "" {
		return nil, status.Error(codes.Unauthenticated, "refresh token required")
	}

	refresh, access, err := t.service.Refresh(ctx, refreshToken)
	if err != nil {
		t.log.Warn("refresh failed", slog.Any("error", err))
		clearRefreshCookie(ctx, t.log)
		return nil, status.Error(
			core_transport_grpc_response.RespondCode(err),
			"failed to refresh token",
		)
	}

	resp := &ssov1.RefreshResponse{AccessToken: access}
	if !fromCookie {
		resp.RefreshToken = refresh
		return resp, nil
	}

	cookieValue := fmt.Sprintf(
		"refresh_token=%s; Path=/; Max-Age=%d; HttpOnly; Secure; SameSite=Lax",
		refresh,
		int(t.service.GetRefreshTokenTTL()/time.Second),
	)

	md := metadata.Pairs("set-cookie", cookieValue)

	if err := grpc.SendHeader(ctx, md); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to set cookie: %v", err)
	}

	return resp, nil
}
