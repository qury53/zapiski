package ai_provider

import (
	"context"

	"github.com/qury53/GC_learn/agent/internal/core/domain"
)

type AIProvider interface {
	GenerateCompletion(ctx context.Context,
		request domain.ChatRequest,
	) (domain.ChatResponse, error)
}
