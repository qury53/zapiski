package agent_service

import (
	chat_postgres_repository "github.com/qury53/GC_learn/agent/internal/repository/chat"
	message_postgres_repository "github.com/qury53/GC_learn/agent/internal/repository/message"
	"github.com/qury53/GC_learn/agent/pkg/ai_provider"
	"github.com/qury53/GC_learn/platform/logger"
)

type AgentService struct {
	log               *logger.Logger
	aiProvider        ai_provider.AIProvider
	messageRepository *message_postgres_repository.MessageRepository
	chatRepository    *chat_postgres_repository.ChatRepository
}

func NewAgentService(provider ai_provider.AIProvider) *AgentService {
	return &AgentService{
		aiProvider: provider,
	}
}
