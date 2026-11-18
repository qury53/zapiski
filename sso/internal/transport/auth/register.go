package auth_transport

import (
	"context"
	"log/slog"
	ssov1 "shared/pkg/proto/sso/v1"

	"github.com/qury53/GC_learn/sso/internal/core/domain"
	core_transport_grpc_response "github.com/qury53/GC_learn/sso/internal/core/transport/grpc"
	"google.golang.org/grpc/status"
)

func (t *AuthTransport) Register(
	ctx context.Context,
	req *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {
	user := domain.NewUserUninitialized(
		req.Username,
		req.Password,
		req.Email,
	)

	if err := user.Validate(); err != nil {
		t.log.Debug("validation error", slog.String("error", err.Error()))
		return nil, status.Error(
			core_transport_grpc_response.RespondCode(err),
			"validation error",
		)
	}

	user_id, verificationStatus, err := t.service.Register(ctx, user)
	if err != nil {
		t.log.Debug("register error", slog.String("error", err.Error()))
		return nil, status.Error(
			core_transport_grpc_response.RespondCode(err),
			"failed to register",
		)
	}

	return &ssov1.RegisterResponse{
		UserId: user_id,
		Status: ssov1.AuthStatus(verificationStatus),
	}, nil
}
