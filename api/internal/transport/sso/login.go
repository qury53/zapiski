package sso

import (
	"encoding/json"
	"net/http"

	core_cookie "github.com/qury53/GC_learn/api/internal/core/http/cookie"
	metadata_carrier "github.com/qury53/GC_learn/api/internal/core/http/metadata"
	core_http_response "github.com/qury53/GC_learn/api/internal/core/http/responce"
	"google.golang.org/grpc/metadata"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	UserID      int64  `json:"user_id"`
}

func (t *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	ctx := r.Context()
	log := t.log
	responseWriter := core_http_response.NewHTTPResponseHandler(log, w)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseWriter.ErrorResponse(err, "invalid body")
		return
	}

	headerMD := metadata.MD{}
	headers := metadata_carrier.NewGRPCCarrier(&headerMD)
	ctx = core_cookie.WithForwardedMetadata(r)

	resp, err := t.service.Login(ctx, req, headers)
	if err != nil {
		responseWriter.ErrorResponse(err, "failed to login")
		return
	}

	core_cookie.WriteSetCookiesFromMetadata(w, headerMD)

	responseWriter.JSONResponse(LoginResponse{
		AccessToken: resp.AccessToken,
		UserID:      resp.UserID,
	}, http.StatusOK)
}
