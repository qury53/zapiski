package sso

import (
	"net/http"

	core_cookie "github.com/qury53/GC_learn/api/internal/core/http/cookie"
	metadata_carrier "github.com/qury53/GC_learn/api/internal/core/http/metadata"
	core_http_response "github.com/qury53/GC_learn/api/internal/core/http/responce"
	"google.golang.org/grpc/metadata"
)

func (t *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := t.log
	responseWriter := core_http_response.NewHTTPResponseHandler(log, w)

	ctx = core_cookie.WithForwardedMetadata(r)

	headerMD := metadata.MD{}
	headers := metadata_carrier.NewGRPCCarrier(&headerMD)
	token := extractRefreshFromCookie(ctx)
	resp, err := t.service.Refresh(ctx, token, headers)
	if err != nil {
		responseWriter.ErrorResponse(err, "failed to refresh")
		return
	}

	core_cookie.WriteSetCookiesFromMetadata(w, headerMD)

	responseWriter.JSONResponse(resp, http.StatusOK)
}
