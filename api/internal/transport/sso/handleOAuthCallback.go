package sso

import (
	"net/http"

	core_cookie "github.com/qury53/GC_learn/api/internal/core/http/cookie"
	metadata_carrier "github.com/qury53/GC_learn/api/internal/core/http/metadata"
	core_http_response "github.com/qury53/GC_learn/api/internal/core/http/responce"
	"google.golang.org/grpc/metadata"
)

type OauthCallbackRequest struct {
	Provider     string `json:"provider"`
	Code         string `json:"code"`
	State        string `json:"state"`
	CodeVerifier string `json:"code_verifier"`
}

type OauthCallbackResponse struct {
	AccessToken string `json:"access_token"`
	UserID      int64  `json:"user_id"`
}

func (t *AuthHandler) HandleOAuthCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := t.log
	rw := core_http_response.NewHTTPResponseHandler(log, w)

	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	provider := r.URL.Query().Get("provider")

	if code == "" || state == "" || provider == "" {
		rw.ErrorResponse(nil, "missing required query parameters")
		return
	}

	ctx = core_cookie.WithForwardedMetadata(r)

	codeVerifier := extractCodeVerifierFromCookie(ctx)
	if codeVerifier == "" {
		rw.ErrorResponse(nil, "missing code_verifier cookie")
		return
	}

	headerMD := metadata.MD{}
	headers := metadata_carrier.NewGRPCCarrier(&headerMD)

	resp, err := t.service.HandleOAuthCallback(ctx, OauthCallbackRequest{
		Provider:     provider,
		Code:         code,
		State:        state,
		CodeVerifier: codeVerifier,
	}, headers)
	if err != nil {
		rw.ErrorResponse(err, "failed to authenticate with provider")
		return
	}

	core_cookie.WriteSetCookiesFromMetadata(w, headerMD)
	rw.JSONResponse(resp, http.StatusOK)
}
