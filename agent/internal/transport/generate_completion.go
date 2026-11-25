package agent_transport

import (
	"context"
	"log/slog"
	agentV1 "shared/pkg/proto/agent/v1"

	"github.com/qury53/GC_learn/agent/internal/domain"
	"google.golang.org/grpc/status"
)

func (t *AgentTransport) ChatCompletion(
	ctx context.Context,
	req *agentV1.ChatRequest,
) (*agentV1.ChatResponse, error) {
	resp, err := t.agentService.GenerateCompletion(ctx, domain.RequestDtoToDomain(req))
	if err != nil {
		t.log.Warn(
			"generate completion error",
			slog.Any("error", err),
		)

		return &agentV1.ChatResponse{}, status.Error(RespondCode(err), err.Error())
	}

	return domain.ResponseDomainToDto(resp), nil
}
