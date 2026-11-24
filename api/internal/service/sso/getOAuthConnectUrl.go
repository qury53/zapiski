package sso

import (
	"context"
	"fmt"

	"github.com/qury53/GC_learn/api/internal/core/ports"
	transport_sso "github.com/qury53/GC_learn/api/internal/transport/sso"
)

func (s *Service) GetOAuthConnectUrl(ctx context.Context, req transport_sso.OauthConnectRequest, headers ports.HeaderCarrier) (transport_sso.OauthConnectResponse, error) {
	resp, err := s.client.GetOAuthConnectUrl(ctx, ports.OauthConnectRequest{
		Provider:            req.Provider,
		RedirectURI:         req.RedirectURI,
		CodeChallenge:       req.CodeChallenge,
		CodeChallengeMethod: req.CodeChallengeMethod,
	}, headers)
	if err != nil {
		return transport_sso.OauthConnectResponse{}, fmt.Errorf("get oauth connect url: %w", err)
	}
	return transport_sso.OauthConnectResponse{AuthorizationURL: resp.AuthorizationURL, OAuthState: resp.OAuthState}, nil
}
