package ai_provider_smartapi

import (
	"fmt"
	"time"

	"github.com/qury53/GC_learn/agent/internal/core/domain"
)

type ChatRequestDto struct {
	MaxTokens   *int         `json:"max_tokens,omitempty"`
	Temperature *float64     `json:"temperature,omitempty"`
	Dto         string       `json:"model"`
	Messages    []MessageDto `json:"messages"`
	Stream      bool         `json:"stream"`
}

type MessageDto struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResponseDto struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
	Created int64    `json:"created"`
}

type Choice struct {
	Message      MessageDto `json:"message"`
	FinishReason string     `json:"finish_reason"`
	Index        int        `json:"index"`
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
		Dto:         req.Model, // поле с тегом `json:"model"`
		Messages:    msgs,
		Stream:      false,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
	}
}

func dtoResponseToDomain(dto ChatResponseDto) (domain.ChatResponse, error) {
	if len(dto.Choices) == 0 {
		return domain.ChatResponse{}, fmt.Errorf("smartapi: response has no choices")
	}

	choice := dto.Choices[0]

	return domain.ChatResponse{
		ThreadID:     dto.ID,
		Model:        dto.Model,
		Content:      choice.Message.Content,
		FinishReason: domain.FinishReason(choice.FinishReason),
		Usage: domain.Usage{
			InputTokens:  dto.Usage.PromptTokens,
			OutputTokens: dto.Usage.CompletionTokens,
			TotalTokens:  dto.Usage.TotalTokens,
			CachedTokens: 0,
		},
		CreatedAt: time.Unix(dto.Created, 0).UTC(),
	}, nil
}
