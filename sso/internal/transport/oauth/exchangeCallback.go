package oauth_transport

import (
	"context"
	"fmt"
	"log/slog"
	ssov1 "shared/pkg/proto/sso/v1"
	"strings"

	core_errors "github.com/qury53/GC_learn/sso/internal/core/errors"
	core_transport_grpc_response "github.com/qury53/GC_learn/sso/internal/core/transport/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// HandleOAuthCallback обменивает authorization code на токены провайдера,
// находит или создаёт локального пользователя и выдаёт свою пару токенов.
// Refresh токен кладётся в HttpOnly cookie, access возвращается в теле ответа.
//
// Returns: *ssov1.LoginResponse, error
func (t *OAuthTransport) HandleOAuthCallback(
	ctx context.Context,
	req *ssov1.OAuthCallbackRequest,
) (*ssov1.LoginResponse, error) {
	provider, err := mapProvider(req.GetProvider())
	if err != nil {
		t.log.Debug("oauth callback: unsupported provider",
			slog.String("provider", req.GetProvider().String()),
		)

		return nil, status.Error(
			core_transport_grpc_response.RespondCode(err),
			"unsupported oauth provider",
		)
	}

	if err := validateCallbackRequest(req); err != nil {
		t.log.Debug("oauth callback: validation error",
			slog.String("provider", string(provider)),
			slog.String("error", err.Error()),
		)

		return nil, status.Error(
			core_transport_grpc_response.RespondCode(err),
			"oauth callback validation error",
		)
	}

	userAgent := getUserAgent(ctx)

	access, refresh, userID, err := t.service.ExchangeCallback(
		ctx,
		provider,
		req.GetCode(),
		req.GetState(),
		req.GetCodeVerifier(),
		&userAgent,
		getClientIP(ctx),
	)
	if err != nil {
		// code, state и verifier в логи не попадают.
		t.logOAuthError("oauth callback failed", provider, err)

		return nil, status.Error(
			core_transport_grpc_response.RespondCode(err),
			"failed to exchange oauth callback",
		)
	}

	if err := t.setRefreshCookie(ctx, refresh); err != nil {
		t.log.Warn("oauth callback: failed to set refresh cookie",
			slog.String("provider", string(provider)),
			slog.Any("error", err),
		)

		return nil, status.Error(codes.Internal, "failed to set refresh cookie")
	}

	return &ssov1.LoginResponse{
		AccessToken: access,
		UserId:      userID,
	}, nil
}

// validateCallbackRequest проверяет обязательные параметры callback.
//
// Returns: error
func validateCallbackRequest(req *ssov1.OAuthCallbackRequest) error {
	if strings.TrimSpace(req.GetCode()) == "" {
		return fmt.Errorf("%w: code is required", core_errors.ErrInvalidArgument)
	}

	if strings.TrimSpace(req.GetState()) == "" {
		return fmt.Errorf("%w: state is required", core_errors.ErrInvalidArgument)
	}

	// code_verifier по RFC 7636 — 43..128 символов base64url.
	verifier := req.GetCodeVerifier()
	if len(verifier) < 43 || len(verifier) > 128 {
		return fmt.Errorf(
			"%w: code_verifier length must be between 43 and 128",
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}
