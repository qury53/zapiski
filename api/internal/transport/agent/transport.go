package agent

import (
	"shared/pkg/logger"
	agentv1 "shared/pkg/proto/agent/v1"
)

type AgentTransport struct {
	log     *logger.Logger
	agentv1 agentv1.AgentServiceClient
}

func NewAgentTransport(log *logger.Logger, agent agentv1.AgentServiceClient) *AgentTransport {
	return &AgentTransport{
		log:     log,
		agentv1: agent,
	}
}
