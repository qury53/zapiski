package user_transport

import (
	"context"
	userv1 "shared/pkg/proto/user/v1"

	"github.com/qury53/GC_learn/platform/auth"
	transport_grpc_response "github.com/qury53/GC_learn/user/internal/core/transport/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (t *UserTransport) GetUserByID(
	ctx context.Context,
	req *userv1.GetUserByIDRequest,
) (*userv1.GetUserByIDResponse, error) {
	callerID, ok := auth.UserIDFrom(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "no identity")
	}

	// Чужой профиль читать нельзя: user_id из запроса должен совпадать с токеном.
	if req.GetUserId() != 0 && req.GetUserId() != callerID {
		return nil, status.Error(codes.PermissionDenied, "forbidden")
	}

	user, err := t.service.GetUserByID(
		ctx,
		callerID,
	)

	if err != nil {
		return nil, status.Error(
			transport_grpc_response.RespondCode(err),
			"GetUserByID transport: failed to get user",
		)
	}

	return &userv1.GetUserByIDResponse{
		User: UserProfileDomainToProto(&user),
	}, nil
}
