package oauth_transport

import (
	"context"
	"fmt"
	"log/slog"
	ssov1 "shared/pkg/proto/sso/v1"
	"strings"

	core_errors "github.com/qury53/GC_learn/sso/internal/core/errors"
	core_transport_grpc_response "github.com/qury53/GC_learn/sso/internal/core/transport/grpc"
	"google.golang.org/grpc/status"
)

// GetOAuthConnectUrl собирает URL страницы согласия провайдера и state,
// который сервис сохранил в кэше для проверки в callback.
//
// Returns: *ssov1.OAuthConnectResponse, error
func (t *OAuthTransport) GetOAuthConnectUrl(
	ctx context.Context,
	req *ssov1.OAuthConnectRequest,
) (*ssov1.OAuthConnectResponse, error) {
	provider, err := mapProvider(req.GetProvider())
	if err != nil {
		t.log.Debug("oauth connect: unsupported provider",
			slog.String("provider", req.GetProvider().String()),
		)

		return nil, status.Error(
			core_transport_grpc_response.RespondCode(err),
			"unsupported oauth provider",
		)
	}

	if err := validateConnectRequest(req); err != nil {
		t.log.Debug("oauth connect: validation error",
			slog.String("provider", string(provider)),
			slog.String("error", err.Error()),
		)

		return nil, status.Error(
			core_transport_grpc_response.RespondCode(err),
			"oauth connect validation error",
		)
	}

	authorizationURL, oauthState, err := t.service.GetConnectURI(
		ctx,
		provider,
		strings.TrimSpace(req.GetRedirectUri()),
		req.GetCodeChallenge(),
	)
	if err != nil {
		t.logOAuthError("oauth connect failed", provider, err)

		return nil, status.Error(
			core_transport_grpc_response.RespondCode(err),
			"failed to build oauth connect url",
		)
	}

	return &ssov1.OAuthConnectResponse{
		AuthorizationUrl: authorizationURL,
		OauthState:       oauthState,
	}, nil
}

// validateConnectRequest проверяет входные данные до вызова сервиса.
// PKCE поддерживаем только в режиме S256: plain не защищает authorization code.
//
// Returns: error
func validateConnectRequest(req *ssov1.OAuthConnectRequest) error {
	if strings.TrimSpace(req.GetRedirectUri()) == "" {
		return fmt.Errorf("%w: redirect_uri is required", core_errors.ErrInvalidArgument)
	}

	if req.GetCodeChallengeMethod() != ssov1.CodeChallengeMethod_CODE_CHALLENGE_METHOD_S256 {
		return fmt.Errorf(
			"%w: only S256 code challenge method is supported",
			core_errors.ErrInvalidArgument,
		)
	}

	// code_challenge для S256 — base64url от sha256, по RFC 7636 это 43..128 символов.
	challenge := req.GetCodeChallenge()
	if len(challenge) < 43 || len(challenge) > 128 {
		return fmt.Errorf(
			"%w: code_challenge length must be between 43 and 128",
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}
