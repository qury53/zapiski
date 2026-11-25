package domain

import (
	agentV1 "shared/pkg/proto/agent/v1"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func RequestDtoToDomain(dto *agentV1.ChatRequest) ChatRequest {
	if dto == nil {
		return ChatRequest{}
	}

	messages := make([]Message, 0, len(dto.GetMessages()))

	for _, msg := range dto.GetMessages() {
		messages = append(messages, Message{
			Role:    RoleFromProto(msg.GetRole()),
			Content: msg.GetContent(),
		})
	}

	result := ChatRequest{
		Model:    dto.GetModel(),
		Messages: messages,
		UserID:   int(dto.GetUserId()),
		ChatMode: ChatModeFromProto(dto.GetMode()),
	}

	if dto.GetChatId() > 0 {
		chatID := int(dto.GetChatId())
		result.ChatID = &chatID
	}

	if dto.MaxTokens != nil {
		maxTokens := int(*dto.MaxTokens)
		result.MaxTokens = &maxTokens
	}

	if dto.Temperature != nil {
		temperature := *dto.Temperature
		result.Temperature = &temperature
	}

	return result
}

// ChatModeFromProto маппит protobuf enum в доменный режим чата.
// UNSPECIFIED означает «как в БД по умолчанию» — full.
func ChatModeFromProto(mode agentV1.ChatMode) ChatMode {
	switch mode {
	case agentV1.ChatMode_CHAT_MODE_SUPPORT:
		return ModeSupport
	case agentV1.ChatMode_CHAT_MODE_FULL:
		return ModeFull
	default:
		return ModeFull
	}
}

func RoleFromProto(role agentV1.Role) Role {
	switch role {
	case agentV1.Role_ROLE_SYSTEM:
		return RoleSystem
	case agentV1.Role_ROLE_USER:
		return RoleUser
	case agentV1.Role_ROLE_ASSISTANT:
		return RoleAssistant
	case agentV1.Role_ROLE_TOOL:
		return RoleTool
	case agentV1.Role_ROLE_UNSPECIFIED:
		fallthrough
	default:
		return ""
	}
}

func ResponseDtoToDomain(dto *agentV1.ChatResponse) ChatResponse {
	if dto == nil {
		return ChatResponse{}
	}

	return ChatResponse{
		ThreadID:     dto.GetThreadId(),
		Model:        dto.GetModel(),
		Content:      dto.GetContent(),
		FinishReason: FinishReasonFromProto(dto.GetFinishReason()),
		Usage:        UsageFromProto(dto.GetUsage()),
		CreatedAt:    dto.GetCreatedAt().AsTime(),
	}
}

func FinishReasonFromProto(reason agentV1.FinishReason) FinishReason {
	switch reason {
	case agentV1.FinishReason_FINISH_REASON_STOP:
		return FinishReasonStop
	case agentV1.FinishReason_FINISH_REASON_LENGTH:
		return FinishReasonLength
	case agentV1.FinishReason_FINISH_REASON_TOOL_CALLS:
		return FinishReasonToolCalls
	case agentV1.FinishReason_FINISH_REASON_CONTENT_FILTER:
		return FinishReasonContentFilter
	case agentV1.FinishReason_FINISH_REASON_UNSPECIFIED:
		fallthrough
	default:
		return ""
	}
}

func UsageFromProto(dto *agentV1.Usage) Usage {
	if dto == nil {
		return Usage{}
	}

	return Usage{
		InputTokens:  int(dto.GetInputTokens()),
		OutputTokens: int(dto.GetOutputTokens()),
		TotalTokens:  int(dto.GetTotalTokens()),
		CachedTokens: int(dto.GetCachedTokens()),
	}
}

func ResponseDomainToDto(resp ChatResponse) *agentV1.ChatResponse {
	return &agentV1.ChatResponse{
		ThreadId:     resp.ThreadID,
		Model:        resp.Model,
		Content:      resp.Content,
		FinishReason: FinishReasonToProto(resp.FinishReason),
		Usage:        UsageToProto(resp.Usage),
		CreatedAt:    timestamppb.New(resp.CreatedAt),
	}
}

func FinishReasonToProto(reason FinishReason) agentV1.FinishReason {
	switch reason {
	case FinishReasonStop:
		return agentV1.FinishReason_FINISH_REASON_STOP
	case FinishReasonLength:
		return agentV1.FinishReason_FINISH_REASON_LENGTH
	case FinishReasonToolCalls:
		return agentV1.FinishReason_FINISH_REASON_TOOL_CALLS
	case FinishReasonContentFilter:
		return agentV1.FinishReason_FINISH_REASON_CONTENT_FILTER
	default:
		return agentV1.FinishReason_FINISH_REASON_UNSPECIFIED
	}
}

func UsageToProto(usage Usage) *agentV1.Usage {
	return &agentV1.Usage{
		InputTokens:  int32(usage.InputTokens),
		OutputTokens: int32(usage.OutputTokens),
		TotalTokens:  int32(usage.TotalTokens),
		CachedTokens: int32(usage.CachedTokens),
	}
}
