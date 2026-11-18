package auth_transport

import (
	"context"
	"log/slog"
	ssov1 "shared/pkg/proto/sso/v1"

	core_transport_grpc_response "github.com/qury53/GC_learn/sso/internal/core/transport/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (t *AuthTransport) Logout(
	ctx context.Context,
	req *ssov1.LogoutRequest,
) (*emptypb.Empty, error) {
	refreshToken := extractRefreshFromCookie(ctx)
	refreshFromCookie := refreshToken != ""
	if !refreshFromCookie {
		refreshToken = req.GetRefreshToken()
	}

	if refreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh token required")
	}

	if err := t.service.Logout(
		ctx,
		refreshToken,
	); err != nil {
		t.log.Warn(
			"failed to logout",
			slog.Any("error", err),
		)
		return nil, status.Error(core_transport_grpc_response.RespondCode(err), "logout failed")
	}

	return nil, nil
}
