package ai_provider_gc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/qury53/GC_learn/agent/internal/domain"
	core_errors "github.com/qury53/GC_learn/agent/internal/errors"
	ai_provider_gc_models "github.com/qury53/GC_learn/agent/pkg/ai_provider/gigachat/models"
)

func (p *GCProvider) GenerateCompletion(
	ctx context.Context,
	request domain.ChatRequest,
) (domain.ChatResponse, error) {
	dtoResponse, err := p.generateCompletion(ctx, ai_provider_gc_models.ToGigaChatRequest(request))
	if err != nil {
		return domain.ChatResponse{}, err
	}

	return ai_provider_gc_models.ToDomainResponse(dtoResponse), nil
}

func (p *GCProvider) generateCompletion(
	ctx context.Context,
	request ai_provider_gc_models.ChatCompletionRequest,
) (ai_provider_gc_models.ChatCompletionResponse, error) {
	log.Printf("generate completion repository")
	token, err := p.refreshAccessToken(ctx)
	if err != nil {
		return ai_provider_gc_models.ChatCompletionResponse{}, fmt.Errorf(
			"token expired: %w: %v",
			core_errors.ErrUnauthorized,
			err,
		)
	}

	url := "https://api.giga.chat/v2/chat/completions"

	payload, err := json.Marshal(request)
	if err != nil {
		return ai_provider_gc_models.ChatCompletionResponse{}, fmt.Errorf("failed to unmarshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(payload))
	if err != nil {
		return ai_provider_gc_models.ChatCompletionResponse{}, fmt.Errorf("failed to make request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Session-ID", "550e8400-e29b-41d4-a716-446655440000")
	response, err := p.client.Do(req)
	if err != nil {
		return ai_provider_gc_models.ChatCompletionResponse{}, fmt.Errorf("failed to do request: %w", err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return ai_provider_gc_models.ChatCompletionResponse{}, fmt.Errorf("failed to read body: %w", err)
	}

	if err := caseStatus(response.StatusCode); err != nil {
		return ai_provider_gc_models.ChatCompletionResponse{}, fmt.Errorf("%w: %s", err, string(body))
	}

	var generateCompletionResponse ai_provider_gc_models.ChatCompletionResponse
	if err := json.Unmarshal(
		body,
		&generateCompletionResponse,
	); err != nil {
		return ai_provider_gc_models.ChatCompletionResponse{}, fmt.Errorf("failed to unmarshal: %w", err)
	}

	return generateCompletionResponse, nil
}

func caseStatus(statusCode int) error {
	switch statusCode {
	case http.StatusOK:
		return nil
	case http.StatusBadRequest:
		return fmt.Errorf(
			"%w",
			core_errors.ErrInvalidArgument,
		)

	case http.StatusUnauthorized:
		return fmt.Errorf(
			"%w",
			core_errors.ErrUnauthorized,
		)

	case http.StatusNotFound:
		return fmt.Errorf(
			"%w",
			core_errors.ErrNotFound,
		)

	case http.StatusUnprocessableEntity:
		return fmt.Errorf(
			"%w",
			core_errors.ErrUnprocessableEntity,
		)

	case http.StatusTooManyRequests:
		return fmt.Errorf(
			"%w",
			core_errors.ErrTooMantRequest,
		)
	}

	return core_errors.ErrUndefinedError
}
