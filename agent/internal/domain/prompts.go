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
	ID        int
	ChatID    int
	Role      Role
	Content   string
	CreatedAt time.Time
}

type ChatRequest struct {
	ChatID      *int
	ChatMode    ChatMode
	Model       string
	Messages    []Message
	MaxTokens   *int
	Temperature *float64
	UserID      int
}

type ChatResponse struct {
	ThreadID     string
	Model        string
	Content      string
	FinishReason FinishReason
	Usage        Usage
	CreatedAt    time.Time
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
