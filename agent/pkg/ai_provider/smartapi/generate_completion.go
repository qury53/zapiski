package ai_provider_smartapi

import (
	"context"

	"github.com/qury53/GC_learn/agent/internal/domain"
)

const (
	baseGenerateCompletionURL = "https://api.smartapi.pro/v1"
)

func (p *SmartAPIProvider) GenerateCompletion(
	ctx context.Context,
	request domain.ChatRequest,
) (domain.ChatResponse, error) {

}
