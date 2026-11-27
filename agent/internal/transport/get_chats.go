package agent_transport

import (
	"context"
	"log/slog"
	v1 "shared/pkg/proto/agent/v1"

	"github.com/qury53/GC_learn/platform/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (t *AgentTransport) GetChats(
	ctx context.Context,
	req *v1.GetChatsRequest,
) (*v1.GetChatsResponse, error) {
	callerID, ok := auth.UserIDFrom(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "no identity")
	}

	if req.GetUserId() != 0 && req.GetUserId() != callerID {
		return nil, status.Error(codes.PermissionDenied, "forbidden")
	}

	chats, err := t.agentService.GetChats(
		ctx,
		int(callerID),
		int64PtrToIntPtr(req.Limit),
		int64PtrToIntPtr(req.Offset),
	)

	if err != nil {
		t.log.Warn(
			"failed to get chats",
			slog.Int64("userID", callerID),
			slog.Any("error", err),
		)

		return nil, status.Error(RespondCode(err), "failed to get chats")
	}

	return ToGetChatsResponse(chats), nil
}
