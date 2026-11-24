package sso

import (
	"encoding/json"
	"net/http"

	core_cookie "github.com/qury53/GC_learn/api/internal/core/http/cookie"
	metadata_carrier "github.com/qury53/GC_learn/api/internal/core/http/metadata"
	core_http_response "github.com/qury53/GC_learn/api/internal/core/http/responce"
	"google.golang.org/grpc/metadata"
)

type OauthConnectRequest struct {
	Provider            string `json:"provider"`
	RedirectURI         string `json:"redirect_uri"`
	CodeChallenge       string `json:"code_challenge"`
	CodeChallengeMethod string `json:"code_challenge_method"`
}

type OauthConnectResponse struct {
	AuthorizationURL string `json:"authorization_url"`
	OAuthState       string `json:"oauth_state"`
}

func (h *AuthHandler) GetOAuthConnectUrl(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := h.log
	rw := core_http_response.NewHTTPResponseHandler(log, w)

	var req OauthConnectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rw.ErrorResponse(err, "invalid body")
		return
	}

	if req.Provider == "" || req.RedirectURI == "" || req.CodeChallenge == "" {
		rw.ErrorResponse(nil, "missing required fields")
		return
	}
	if req.CodeChallengeMethod == "" {
		req.CodeChallengeMethod = "S256"
	}

	ctx = core_cookie.WithForwardedMetadata(r)
	headerMD := metadata.MD{}
	headers := metadata_carrier.NewGRPCCarrier(&headerMD)

	resp, err := h.service.GetOAuthConnectUrl(ctx, req, headers)
	if err != nil {
		rw.ErrorResponse(err, "failed to get oauth connect url")
		return
	}

	core_cookie.WriteSetCookiesFromMetadata(w, headerMD)
	rw.JSONResponse(resp, http.StatusOK)
}
