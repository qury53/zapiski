package ai_provider_gc_models

import "encoding/json"

// ============================================================================
// Request
// ============================================================================

type ChatCompletionRequest struct {
	ThreadID      *string        `json:"thread_id,omitempty"`
	CreatedAt     *int64         `json:"created_at,omitempty"`
	ModelOptions  *ModelOptions  `json:"model_options,omitempty"`
	RankerOptions *RankerOptions `json:"ranker_options,omitempty"`
	UserInfo      *UserInfo      `json:"user_info,omitempty"`
	ToolConfig    *ToolConfig    `json:"tool_config,omitempty"`
	Model         string         `json:"model"`
	Messages      []Message      `json:"messages"`
	Tools         []Tool         `json:"tools,omitempty"`
	Stream        bool           `json:"stream,omitempty"`
	DisableFilter bool           `json:"disable_filter,omitempty"`
}

type ModelOptions struct {
	Temperature         *float64        `json:"temperature,omitempty"`
	TopP                *float64        `json:"top_p,omitempty"`
	MaxTokens           *int32          `json:"max_tokens,omitempty"`
	RepetitionPenalty   *float64        `json:"repetition_penalty,omitempty"`
	UpdateInterval      *int32          `json:"update_interval,omitempty"`
	TopLogprobs         *int32          `json:"top_logprobs,omitempty"`
	ResponseFormat      *ResponseFormat `json:"response_format,omitempty"`
	ParallelToolCalls   *bool           `json:"parallel_tool_calls,omitempty"`
	UnnormalizedHistory bool            `json:"unnormalized_history,omitempty"`
}

type ResponseFormat struct {
	Type   string          `json:"type"`
	Schema json.RawMessage `json:"schema,omitempty"`
	Strict bool            `json:"strict,omitempty"`
}

type RankerOptions struct {
	TopN            *int32  `json:"top_n,omitempty"`
	EmbeddingsModel *string `json:"embeddings_model,omitempty"`
	Enabled         bool    `json:"enabled"`
}

type UserInfo struct {
	Timezone string `json:"timezone"`
}

type ToolConfig struct {
	ToolName     *string `json:"tool_name,omitempty"`
	FunctionName *string `json:"function_name,omitempty"`
	Mode         string  `json:"mode"`
}

type Tool struct {
	// В JSON одновременно должен быть задан только один вариант.
	Functions       *Functions       `json:"functions,omitempty"`
	ImageGenerate   *ImageGenerate   `json:"image_generate,omitempty"`
	Model3DGenerate *Model3DGenerate `json:"model_3d_generate,omitempty"`
}

type Functions struct {
	Specifications []FunctionSpecification `json:"specifications"`
}

type FunctionSpecification struct {
	Name             string           `json:"name"`
	Description      string           `json:"description"`
	Parameters       json.RawMessage  `json:"parameters,omitempty"`
	FewShotExamples  []FewShotExample `json:"few_shot_examples,omitempty"`
	ReturnParameters json.RawMessage  `json:"return_parameters,omitempty"`
}

type FewShotExample struct {
	Request string          `json:"request"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type ImageGenerate struct{}

type Model3DGenerate struct{}

// ============================================================================
// Shared messages
// ============================================================================

type Role string

const (
	RoleUser      Role = "user"
	RoleSystem    Role = "system"
	RoleAssistant Role = "assistant"
	RoleTool      Role = "tool"
)

type Message struct {
	MessageID   *string       `json:"message_id,omitempty"`
	Role        Role          `json:"role"`
	ToolStateID *string       `json:"tool_state_id,omitempty"`
	Content     []ContentPart `json:"content"`
}

type ContentPart struct {
	Text           *string         `json:"text,omitempty"`
	InlineData     *InlineData     `json:"inline_data,omitempty"`
	Files          []FileRef       `json:"files,omitempty"`
	FunctionCall   *FunctionCall   `json:"function_call,omitempty"`
	FunctionResult *FunctionResult `json:"function_result,omitempty"`
	ToolExecution  *ToolExecution  `json:"tool_execution,omitempty"`
	Logprobs       []Logprob       `json:"logprobs,omitempty"`
}

type FileRef struct {
	ID     string `json:"id"`
	Target string `json:"target,omitempty"`
	MIME   string `json:"mime,omitempty"`
}

type FunctionCall struct {
	ID   string          `json:"id,omitempty"`
	Name string          `json:"name"`
	Args json.RawMessage `json:"arguments,omitempty"`
}

type FunctionResult struct {
	Name   string `json:"name"`
	Result string `json:"result"`
}

type ToolExecution struct {
	SecondsLeft *int32              `json:"seconds_left,omitempty"`
	Name        string              `json:"name"`
	Status      ToolExecutionStatus `json:"status"`
	Censored    bool                `json:"censored,omitempty"`
}

type ToolExecutionStatus string

const (
	ToolExecutionStatusSuccess ToolExecutionStatus = "success"
	ToolExecutionStatusFailed  ToolExecutionStatus = "failed"
	ToolExecutionStatusPending ToolExecutionStatus = "pending"
)

type Logprob struct {
	Chosen *TokenLogprob  `json:"chosen,omitempty"`
	Top    []TokenLogprob `json:"top,omitempty"`
}

type TokenLogprob struct {
	Token   string  `json:"token"`
	TokenID int32   `json:"token_id"`
	Logprob float64 `json:"logprob"`
}

type InlineData struct {
	Sources map[string]Source        `json:"sources,omitempty"`
	Images  []map[string]interface{} `json:"images,omitempty"`
}

type Source struct {
	Title string `json:"title,omitempty"`
	URL   string `json:"url,omitempty"`
}

// ============================================================================
// Response
// ============================================================================

type ChatCompletionResponse struct {
	Model        string       `json:"model"`
	ThreadID     string       `json:"thread_id"`
	FinishReason FinishReason `json:"finish_reason"`
	Messages     []Message    `json:"messages"`
	CreatedAt    int64        `json:"created_at"`
	Usage        Usage        `json:"usage"`
}

type FinishReason string

const (
	FinishReasonStop          FinishReason = "stop"
	FinishReasonLength        FinishReason = "length"
	FinishReasonToolCalls     FinishReason = "tool_calls"
	FinishReasonContentFilter FinishReason = "content_filter"
)

type Usage struct {
	InputTokens        int32              `json:"input_tokens"`
	InputTokensDetails InputTokensDetails `json:"input_tokens_details"`
	OutputTokens       int32              `json:"output_tokens"`
	TotalTokens        int32              `json:"total_tokens"`
}

type InputTokensDetails struct {
	CachedTokens int32 `json:"cached_tokens"`
}
