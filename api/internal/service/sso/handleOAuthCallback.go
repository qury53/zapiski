package sso

import (
	"context"
	"fmt"

	"github.com/qury53/GC_learn/api/internal/core/ports"
	transport_sso "github.com/qury53/GC_learn/api/internal/transport/sso"
)

func (s *Service) HandleOAuthCallback(ctx context.Context, req transport_sso.OauthCallbackRequest, headers ports.HeaderCarrier) (transport_sso.OauthCallbackResponse, error) {
	resp, err := s.client.HandleOAuthCallback(ctx, ports.OauthCallbackRequest{
		Provider:     req.Provider,
		Code:         req.Code,
		State:        req.State,
		CodeVerifier: req.CodeVerifier,
	}, headers)
	if err != nil {
		return transport_sso.OauthCallbackResponse{}, fmt.Errorf("handle oauth callback: %w", err)
	}
	return transport_sso.OauthCallbackResponse{AccessToken: resp.AccessToken, UserID: resp.UserID}, nil
}
