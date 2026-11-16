package sso

import (
	"encoding/json"
	"fmt"
	"net/http"
	ssov1 "shared/pkg/proto/sso/v1"

	core_http_response "github.com/qury53/GC_learn/api/internal/core/http/responce"
	core_logger "github.com/qury53/GC_learn/api/internal/core/logger"
)

type RefreshRequest struct {
	RefreshToken string
}

func (t *SSOTransport) Refresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest

	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseWriter := core_http_response.NewHTTPResponseHandler(log, w)
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		responseWriter.ErrorResponse(err, "refresh_token required")
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseWriter.ErrorResponse(err, "invalid body")
		return
	}

	resp, err := t.ssov1.Refresh(ctx, &ssov1.RefreshRequest{
		RefreshToken: fmt.Sprintf("%s", cookie.Value),
	})

	if err != nil {
		responseWriter.ErrorResponse(err, "failed to refresh")
		return
	}

	responseWriter.JSONResponse(resp, http.StatusOK)
}
