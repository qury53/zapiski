package user_transport

import (
	"context"
	userv1 "shared/pkg/proto/user/v1"

	transport_grpc_response "github.com/qury53/GC_learn/user/internal/core/transport/grpc"
	"google.golang.org/grpc/status"
)

func (t *UserTransport) GetUsers(
	ctx context.Context,
	req *userv1.GetUsersRequest,
) (*userv1.GetUsersResponse, error) {
	users, err := t.service.GetUsers(ctx, req.Limit, req.Offset)
	if err != nil {
		return &userv1.GetUsersResponse{Users: []*userv1.UserProfile{}}, status.Error(
			transport_grpc_response.RespondCode(err),
			"GetUsers: failed to get users",
		)
	}

	return UserProfilesDomainToProto(users), nil
}
