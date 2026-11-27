package ai_provider_gc_models

import (
	"time"

	"github.com/qury53/GC_learn/agent/internal/domain"
)

func ToGigaChatRequest(req domain.ChatRequest) ChatCompletionRequest {
	messages := make([]Message, 0, len(req.Messages))

	for _, m := range req.Messages {
		messages = append(messages, Message{
			Role: Role(m.Role),
			Content: []ContentPart{
				{Text: &m.Content},
			},
		})
	}

	apiReq := ChatCompletionRequest{
		Model:    req.Model,
		Messages: messages,
	}

	if req.MaxTokens != nil {
		apiReq.ModelOptions = &ModelOptions{
			MaxTokens: ptr(int32(*req.MaxTokens)),
		}
	}

	if req.Temperature != nil {
		if apiReq.ModelOptions == nil {
			apiReq.ModelOptions = &ModelOptions{}
		}
		apiReq.ModelOptions.Temperature = req.Temperature
	}

	return apiReq
}

func ToDomainResponse(resp ChatCompletionResponse) domain.ChatResponse {
	var text string

	for _, msg := range resp.Messages {
		for _, part := range msg.Content {
			if part.Text != nil {
				text += *part.Text
			}
		}
	}

	usage := domain.Usage{
		InputTokens:  int(resp.Usage.InputTokens),
		OutputTokens: int(resp.Usage.OutputTokens),
		TotalTokens:  int(resp.Usage.TotalTokens),
		CachedTokens: int(resp.Usage.InputTokensDetails.CachedTokens),
	}

	return domain.ChatResponse{
		ThreadID:     resp.ThreadID,
		Model:        resp.Model,
		Content:      text,
		FinishReason: domain.FinishReason(resp.FinishReason),
		Usage:        usage,
		CreatedAt:    time.Unix(resp.CreatedAt, 0),
	}
}

func ptr[T any](v T) *T {
	return &v
}
