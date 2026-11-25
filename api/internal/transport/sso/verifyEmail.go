package sso

import (
	"encoding/json"
	"net/http"

	core_cookie "github.com/qury53/GC_learn/api/internal/core/http/cookie"
	metadata_carrier "github.com/qury53/GC_learn/api/internal/core/http/metadata"
	core_http_response "github.com/qury53/GC_learn/api/internal/core/http/responce"
	"google.golang.org/grpc/metadata"
)

type VerifyEmailRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type AuthStatus string

const (
	AuthStatusActive              AuthStatus = "AUTH_STATUS_ACTIVE"
	AuthStatusPendingVerification AuthStatus = "AUTH_STATUS_PENDING_VERIFICATION"
	AuthStatusBlocked             AuthStatus = "AUTH_STATUS_BLOCKED"
)

type VerifyEmailResponse struct {
	Status AuthStatus `json:"status"`
}

func (t *AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	var req VerifyEmailRequest

	ctx := r.Context()
	log := t.log
	responseWriter := core_http_response.NewHTTPResponseHandler(log, w)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseWriter.ErrorResponse(err, "invalid body")
		return
	}

	ctx = core_cookie.WithForwardedMetadata(r)
	headerMD := metadata.MD{}
	headers := metadata_carrier.NewGRPCCarrier(&headerMD)

	resp, err := t.service.VerifyEmail(ctx, req, headers)
	if err != nil {
		responseWriter.ErrorResponse(err, "failed to verify email")
		return
	}

	core_cookie.WriteSetCookiesFromMetadata(w, headerMD)
	responseWriter.JSONResponse(resp, http.StatusOK)
}
