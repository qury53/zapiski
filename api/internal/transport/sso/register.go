package sso

import (
	"encoding/json"
	"net/http"

	core_http_response "github.com/qury53/GC_learn/api/internal/core/http/responce"
	"github.com/qury53/GC_learn/api/internal/core/ports"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type RegisterResponse struct {
	UserID int64  `json:"user_id"`
	Status string `json:"status"`
}

func (t *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	ctx := r.Context()
	log := t.log
	responseWriter := core_http_response.NewHTTPResponseHandler(log, w)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseWriter.ErrorResponse(err, "invalid body")
		return
	}

	resp, err := t.service.Register(ctx, req)
	if err != nil {
		responseWriter.ErrorResponse(err, "failed to register user")
		return
	}

	responseWriter.JSONResponse(RegisterResponse{
		UserID: resp.UserID,
		Status: string(resp.Status),
	}, http.StatusOK)
}

func toPortRegisterRequest(req RegisterRequest) ports.RegisterRequest {
	return ports.RegisterRequest{Username: req.Username, Password: req.Password, Email: req.Email}
}

func fromPortRegisterResponse(resp ports.RegisterResponse) RegisterResponse {
	return RegisterResponse{UserID: resp.UserID, Status: string(resp.Status)}
}
