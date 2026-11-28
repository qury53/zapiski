package agent_transport

import (
	"context"
	agentV1 "shared/pkg/proto/agent/v1"

	"github.com/qury53/GC_learn/agent/internal/core/domain"
	"github.com/qury53/GC_learn/platform/logger"
)

type AgentTransport struct {
	agentV1.UnimplementedAgentServiceServer
	log          *logger.Logger
	agentService AgentService
}

type AgentService interface {
	GenerateCompletion(
		ctx context.Context,
		request domain.ChatRequest,
	) (domain.ChatResponse, error)

	GetChats(
		ctx context.Context,
		user_id int,
		limit *int,
		offset *int,
	) ([]domain.Chat, error)

	GetMessages(
		ctx context.Context,
		chat_id int,
		user_id int,
		limit *int,
		offset *int,
	) ([]domain.Message, error)
}

func NewAgentTransport(
	service AgentService,
	log *logger.Logger,
) *AgentTransport {
	return &AgentTransport{
		log:          log,
		agentService: service,
	}
}
