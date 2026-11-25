package agent_transport

import (
	"context"
	"log/slog"
	v1 "shared/pkg/proto/agent/v1"

	"google.golang.org/grpc/status"
)

func (t *AgentTransport) GetChats(
	ctx context.Context,
	req *v1.GetChatsRequest,
) (*v1.GetChatsResponse, error) {
	chats, err := t.agentService.GetChats(
		ctx,
		int(req.GetUserId()),
		int64PtrToIntPtr(req.Limit),
		int64PtrToIntPtr(req.Offset),
	)

	if err != nil {
		t.log.Warn(
			"failed to get chats",
			slog.Int64("userID", req.GetUserId()),
			slog.Any("error", err),
		)

		return nil, status.Error(RespondCode(err), "failed to get chats")
	}

	return ToGetChatsResponse(chats), nil
}
