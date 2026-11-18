package sso

import (
	"net/http"
	ssov1 "shared/pkg/proto/sso/v1"

	core_http_response "github.com/qury53/GC_learn/api/internal/core/http/responce"
	core_logger "github.com/qury53/GC_learn/api/internal/core/logger"
)

func parseOAuthProvider(providerStr string) ssov1.OAuthProvider {
	switch providerStr {
	case "google":
		return ssov1.OAuthProvider_OAUTH_PROVIDER_GOOGLE
	case "github":
		return ssov1.OAuthProvider_OAUTH_PROVIDER_GITHUB
	case "yandex":
		return ssov1.OAuthProvider_OAUTH_PROVIDER_YANDEX
	default:
		return ssov1.OAuthProvider_OAUTH_PROVIDER_UNSPECIFIED
	}
}

func (t *SSOTransport) HandleOAuthCallback(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseWriter := core_http_response.NewHTTPResponseHandler(log, w)
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	providerFromQuery := r.URL.Query().Get("provider")
	provider := parseOAuthProvider(providerFromQuery)

	resp, err := t.ssov1.HandleOAuthCallback(ctx, &ssov1.OAuthCallbackRequest{
		Provider: provider,
		Code:     code,
		State:    state,
	})
	if err != nil {
		responseWriter.ErrorResponse(err, "failed to authenticate with provider")
		return
	}

	responseWriter.JSONResponse(resp, http.StatusOK)
}
