package agent

import (
	agentv1 "shared/pkg/proto/agent/v1"

	"github.com/qury53/GC_learn/platform/logger"
)

type AgentTransport struct {
	log   *logger.Logger
	agent agentv1.AgentServiceClient
}

func NewAgentTransport(log *logger.Logger, agent agentv1.AgentServiceClient) *AgentTransport {
	return &AgentTransport{
		log:   log,
		agent: agent,
	}
}
