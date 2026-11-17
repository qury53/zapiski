package auth_transport

import (
	"context"
	ssov1 "shared/pkg/proto/sso/v1"

	core_transport_grpc_response "github.com/qury53/GC_learn/sso/internal/core/transport/grpc"
	"google.golang.org/grpc/status"
)

func (t *AuthTransport) VerifyEmail(
	ctx context.Context,
	req *ssov1.VerifyEmailRequest,
) (*ssov1.VerifyEmailResponse, error) {
	verificationStatus, err := t.service.VerifyEmail(
		ctx,
		req.GetEmail(),
		req.GetCode(),
	)

	if err != nil {
		return nil, status.Error(core_transport_grpc_response.RespondCode(err), "invalid code")
	}

	return &ssov1.VerifyEmailResponse{
		Status: ssov1.AuthStatus(verificationStatus),
	}, nil
}
