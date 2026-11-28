package agent_transport

import (
	"context"
	"log/slog"
	v1 "shared/pkg/proto/agent/v1"

	core_errors "github.com/qury53/GC_learn/agent/internal/core/errors"
	"github.com/qury53/GC_learn/platform/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (t *AgentTransport) GetMessages(
	ctx context.Context,
	req *v1.GetMessagesRequest,
) (*v1.GetMessagesResponse, error) {
	callerID, ok := auth.UserIDFrom(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "no identity")
	}

	if req.GetUserId() != 0 && req.GetUserId() != callerID {
		return nil, status.Error(codes.PermissionDenied, "forbidden")
	}

	messages, err := t.agentService.GetMessages(
		ctx,
		int(req.GetChatId()),
		int(callerID),
		int64PtrToIntPtr(req.Limit),
		int64PtrToIntPtr(req.Offset),
	)

	if err != nil {
		t.log.Warn(
			"failed to get messages",
			slog.Int64("chatID", req.ChatId),
			slog.Any("error", err),
		)

		return nil, status.Error(core_errors.RespondCode(err), "failed to get messages")
	}

	return ToGetMessagesResponse(messages), nil
}
