package sso

import (
	"encoding/json"
	"net/http"
	ssov1 "shared/pkg/proto/sso/v1"

	core_http_response "github.com/qury53/GC_learn/api/internal/core/http/responce"
	core_logger "github.com/qury53/GC_learn/api/internal/core/logger"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (t *SSOTransport) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseWriter := core_http_response.NewHTTPResponseHandler(log, w)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseWriter.ErrorResponse(err, "invalid body")
		return
	}

	resp, err := t.ssov1.Register(ctx, &ssov1.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	})
	if err != nil {
		responseWriter.ErrorResponse(err, "failed to register user")
		return
	}

	responseWriter.JSONResponse(resp, http.StatusOK)
}
