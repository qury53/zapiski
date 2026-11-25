package adapters_sso

import (
	"context"
	ssov1 "shared/pkg/proto/sso/v1"

	"github.com/qury53/GC_learn/api/internal/core/ports"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (c *GRPCClient) HandleOAuthCallback(ctx context.Context, req ports.OauthCallbackRequest, headers ports.HeaderCarrier) (ports.OauthCallbackResponse, error) {
	provider, err := oauthProviderToPB(req.Provider)
	if err != nil {
		return ports.OauthCallbackResponse{}, err
	}
	var headerMD, trailerMD metadata.MD
	resp, err := c.client.HandleOAuthCallback(ctx, &ssov1.OAuthCallbackRequest{
		Provider:     provider,
		Code:         req.Code,
		State:        req.State,
		CodeVerifier: req.CodeVerifier,
	}, grpc.Header(&headerMD), grpc.Trailer(&trailerMD))
	copySetCookies(headers, headerMD, trailerMD)
	if err != nil {
		return ports.OauthCallbackResponse{}, err
	}
	return ports.OauthCallbackResponse{AccessToken: resp.GetAccessToken(), UserID: resp.GetUserId()}, nil
}
