package agent_transport

import (
	"context"
	"log/slog"
	agentV1 "shared/pkg/proto/agent/v1"

	"github.com/qury53/GC_learn/agent/internal/core/domain"
	core_errors "github.com/qury53/GC_learn/agent/internal/core/errors"
	"github.com/qury53/GC_learn/platform/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (t *AgentTransport) ChatCompletion(
	ctx context.Context,
	req *agentV1.ChatRequest,
) (*agentV1.ChatResponse, error) {
	callerID, ok := auth.UserIDFrom(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "no identity")
	}

	if req.GetUserId() != 0 && req.GetUserId() != callerID {
		return nil, status.Error(codes.PermissionDenied, "forbidden")
	}

	// Identity берём из токена, а не из запроса.
	req.UserId = callerID

	resp, err := t.agentService.GenerateCompletion(ctx, domain.RequestDtoToDomain(req))
	if err != nil {
		t.log.Warn(
			"generate completion error",
			slog.Any("error", err),
		)

		return &agentV1.ChatResponse{}, status.Error(core_errors.RespondCode(err), err.Error())
	}

	return domain.ResponseDomainToDto(resp), nil
}
