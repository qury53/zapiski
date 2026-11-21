package ai_provider_smartapi

import (
	"time"

	"github.com/qury53/GC_learn/agent/internal/domain"
)

type ChatRequestDto struct {
	Dto      string       `json:"model"`
	Messages []MessageDto `json:"messages"`
	Stream   bool         `json:"stream"`
}

type MessageDto struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResponseDto struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Index        int        `json:"index"`
	Message      MessageDto `json:"message"`
	FinishReason string     `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func domainReqToDto(req domain.ChatRequest) ChatRequestDto {
	msgs := make([]MessageDto, len(req.Messages))
	for i, m := range req.Messages {
		msgs[i] = MessageDto{
			Role:    string(m.Role),
			Content: m.Content,
		}
	}

	return ChatRequestDto{
		Dto:      req.Model, // поле с тегом `json:"model"`
		Messages: msgs,
		Stream:   false,
	}
}

func dtoResponseToDomain(dto ChatResponseDto) domain.ChatResponse {
	var content string
	var finishReason domain.FinishReason

	if len(dto.Choices) > 0 {
		content = dto.Choices[0].Message.Content
		finishReason = domain.FinishReason(dto.Choices[0].FinishReason)
	}

	return domain.ChatResponse{
		ThreadID:     dto.ID,
		Model:        dto.Model,
		Content:      content,
		FinishReason: finishReason,
		Usage: domain.Usage{
			InputTokens:  dto.Usage.PromptTokens,
			OutputTokens: dto.Usage.CompletionTokens,
			TotalTokens:  dto.Usage.TotalTokens,
			CachedTokens: 0,
		},
		CreatedAt: time.Unix(dto.Created, 0).UTC(),
	}
}
