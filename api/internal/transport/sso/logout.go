package sso

import (
	"encoding/json"
	"net/http"
	ssov1 "shared/pkg/proto/sso/v1"

	core_http_response "github.com/qury53/GC_learn/api/internal/core/http/responce"
	core_logger "github.com/qury53/GC_learn/api/internal/core/logger"
)

type LogoutRequest struct {
	RefreshToken string
}

func (t *SSOTransport) Logout(w http.ResponseWriter, r *http.Request) {
	var req LogoutRequest

	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	// TODO: ЕСЛИ ЧТО ПРИДЕЛАТЬ ПАРСЕР КУКИ ПО refresh_token
	responseWriter := core_http_response.NewHTTPResponseHandler(log, w)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseWriter.ErrorResponse(err, "invalid body")
		return
	}

	resp, err := t.ssov1.Logout(ctx, &ssov1.LogoutRequest{
		RefreshToken: req.RefreshToken,
	})

	if err != nil {
		responseWriter.ErrorResponse(err, "failed to logout")
		return
	}

	responseWriter.JSONResponse(resp, http.StatusOK)
}
