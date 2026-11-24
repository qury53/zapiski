package adapters_sso

import (
	"context"
	ssov1 "shared/pkg/proto/sso/v1"

	"github.com/qury53/GC_learn/api/internal/core/ports"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (c *GRPCClient) GetOAuthConnectUrl(ctx context.Context, req ports.OauthConnectRequest, headers ports.HeaderCarrier) (ports.OauthConnectResponse, error) {
	provider, err := oauthProviderToPB(req.Provider)
	if err != nil {
		return ports.OauthConnectResponse{}, err
	}
	method, err := codeChallengeMethodToPB(req.CodeChallengeMethod)
	if err != nil {
		return ports.OauthConnectResponse{}, err
	}
	var headerMD, trailerMD metadata.MD
	resp, err := c.client.GetOAuthConnectUrl(ctx, &ssov1.OAuthConnectRequest{
		Provider:            provider,
		RedirectUri:         req.RedirectURI,
		CodeChallenge:       req.CodeChallenge,
		CodeChallengeMethod: method,
	}, grpc.Header(&headerMD), grpc.Trailer(&trailerMD))
	copySetCookies(headers, headerMD, trailerMD)
	if err != nil {
		return ports.OauthConnectResponse{}, err
	}
	return ports.OauthConnectResponse{AuthorizationURL: resp.GetAuthorizationUrl(), OAuthState: resp.GetOauthState()}, nil
}
