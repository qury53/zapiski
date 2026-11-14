package auth_transport

import (
	"context"
	ssov1 "shared/pkg/proto/sso/v1"
)

func (t *AuthTransport) VerifyEmail(
	ctx context.Context,
	req *ssov1.VerifyEmailRequest,
) (*ssov1.VerifyEmailResponse, error) {
	status := t.service.VerifyEmail(
		ctx,
		req.GetEmail(),
		req.GetCode(),
	)

	return &ssov1.VerifyEmailResponse{
		Status: ssov1.AuthStatus(status),
	}, nil
}
