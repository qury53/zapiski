package ai_provider_smartapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/qury53/GC_learn/agent/internal/core/domain"
	core_errors "github.com/qury53/GC_learn/agent/internal/core/errors"
)

const (
	baseGenerateCompletionURL = "https://api.smartapi.pro/v1/chat/completions"
)

func (p *SmartAPIProvider) GenerateCompletion(
	ctx context.Context,
	request domain.ChatRequest,
) (domain.ChatResponse, error) {
	dtoReq := domainReqToDto(request)

	p.log.Debug(
		"starting completion generation",
		slog.String("model", dtoReq.Dto),
		slog.Int("messages", len(dtoReq.Messages)),
	)

	payload, err := json.Marshal(dtoReq)
	if err != nil {
		return domain.ChatResponse{}, fmt.Errorf("marshaling failed: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		baseGenerateCompletionURL,
		bytes.NewReader(payload),
	)
	if err != nil {
		return domain.ChatResponse{}, fmt.Errorf(
			"required initialization failed: %w",
			err,
		)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.secretKey)

	response, err := p.client.Do(req)
	if err != nil {
		return domain.ChatResponse{}, fmt.Errorf(
			"failed to do requrest: %w",
			err,
		)
	}

	defer func() { _ = response.Body.Close() }()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return domain.ChatResponse{}, fmt.Errorf("read body error : %w", err)
	}

	if err := caseStatus(response.StatusCode); err != nil {
		p.log.Warn(
			"smartapi request failed",
			slog.Int("status", response.StatusCode),
			slog.String("body", truncate(string(body), 512)),
		)

		return domain.ChatResponse{}, err
	}

	var chatResponseDto ChatResponseDto
	if err := json.Unmarshal(body, &chatResponseDto); err != nil {
		return domain.ChatResponse{}, fmt.Errorf("failed to unmarshal: %w", err)
	}

	return dtoResponseToDomain(chatResponseDto)
}

func caseStatus(statusCode int) error {
	switch statusCode {
	case http.StatusOK:
		return nil
	case http.StatusBadRequest, http.StatusNotFound, http.StatusUnprocessableEntity:
		return core_errors.ErrInvalidArgument
	case http.StatusUnauthorized, http.StatusForbidden:
		return core_errors.ErrUnauthorized
	case http.StatusPaymentRequired:
		return core_errors.InsufficientBalance
	case http.StatusTooManyRequests:
		return core_errors.ErrTooMantRequest
	default:
		return core_errors.ErrInternalServerError
	}
}

func truncate(value string, limit int) string {
	if len(value) <= limit {
		return value
	}

	return value[:limit] + "..."
}
