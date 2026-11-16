package sso

import (
	"encoding/json"
	"net/http"
	ssov1 "shared/pkg/proto/sso/v1"

	core_http_response "github.com/qury53/GC_learn/api/internal/core/http/responce"
	core_logger "github.com/qury53/GC_learn/api/internal/core/logger"
)

type VerifyEmailRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

func (t *SSOTransport) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	var req VerifyEmailRequest

	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseWriter := core_http_response.NewHTTPResponseHandler(log, w)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseWriter.ErrorResponse(err, "invalid body")
		return
	}

	resp, err := t.ssov1.VerifyEmail(ctx, &ssov1.VerifyEmailRequest{
		Email: req.Email,
		Code:  req.Code,
	})
	if err != nil {
		responseWriter.ErrorResponse(err, "failed to verify email")
		return
	}

	responseWriter.JSONResponse(resp, http.StatusOK)
}
