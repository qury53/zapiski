package agent_transport

import (
	"context"
	"log/slog"
	v1 "shared/pkg/proto/agent/v1"

	"google.golang.org/grpc/status"
)

func (t *AgentTransport) GetMessages(
	ctx context.Context,
	req *v1.GetMessagesRequest,
) (*v1.GetMessagesResponse, error) {
	messages, err := t.agentService.GetMessages(
		ctx,
		int(req.GetChatId()),
		int64PtrToIntPtr(req.Limit),
		int64PtrToIntPtr(req.Offset),
	)

	if err != nil {
		t.log.Warn(
			"failed to get messages",
			slog.Int64("chatID", req.ChatId),
			slog.Any("error", err),
		)

		return nil, status.Error(RespondCode(err), "failed to get messages")
	}

	return ToGetMessagesResponse(messages), nil
}
