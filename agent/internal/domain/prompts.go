package domain

import "time"

type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleTool      Role = "tool"
)

type Message struct {
	CreatedAt time.Time
	Role      Role
	Content   string
	ID        int
	ChatID    int
}

type ChatRequest struct {
	ChatID      *int
	MaxTokens   *int
	Temperature *float64
	Model       string
	Messages    []Message
	ChatMode    ChatMode
	UserID      int
}

type ChatResponse struct {
	CreatedAt    time.Time
	ThreadID     string
	Model        string
	Content      string
	FinishReason FinishReason
	Usage        Usage
}

type FinishReason string

const (
	FinishReasonStop          FinishReason = "stop"
	FinishReasonLength        FinishReason = "length"
	FinishReasonToolCalls     FinishReason = "tool_calls"
	FinishReasonContentFilter FinishReason = "content_filter"
)

type Usage struct {
	InputTokens  int
	OutputTokens int
	TotalTokens  int
	CachedTokens int
}
