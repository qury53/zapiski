package ai_provider_smartapi

import "github.com/qury53/GC_learn/agent/internal/domain"

type ChatRequestDto struct {
	Dto      string       `json:"model"`
	Messages []MessageDto `json:"messages"`
	Stream   bool         `json:"stream"`
}

type MessageDto struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResponseDto struct{}

func domainReqToDto(domain.ChatRequest) ChatRequestDto {
	//some convertion
	return ChatRequestDto{}
}

func dtoResponseToDomain(dto ChatResponseDto) domain.ChatResponse
