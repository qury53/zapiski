package ai_provider_smartapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/qury53/GC_learn/agent/internal/domain"
)

const (
	baseGenerateCompletionURL = "https://api.smartapi.pro/v1"
)

func (p *SmartAPIProvider) GenerateCompletion(
	ctx context.Context,
	request domain.ChatRequest,
) (domain.ChatResponse, error) {
	domainReq := domainReqToDto(request)
	payload, err := json.Marshal(domainReq)

	if err != nil {
		return domain.ChatResponse{}, fmt.Errorf("marshaling failed: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		baseGenerateCompletionURL,
		bytes.NewReader(payload),
	)

	req.Header.Set("Authorization", "Bearer "+p.secretKey)
	if err != nil {
		return domain.ChatResponse{}, fmt.Errorf(
			"requred initialization failed: %w",
			err,
		)
	}

	response, err := p.client.Do(req)
	if err != nil {
		return domain.ChatResponse{}, fmt.Errorf(
			"failed to do requrest: %w",
			err,
		)
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return domain.ChatResponse{}, fmt.Errorf("read body error : %w", err)
	}

	if err := caseStatus(response.StatusCode); err != nil {
		return domain.ChatResponse{}, err
	}

	var chatResponseDto ChatResponseDto
	if err := json.Unmarshal(body, &chatResponseDto); err != nil {
		return domain.ChatResponse{}, fmt.Errorf("failed to unmarshal: %w", err)
	}

	return dtoResponseToDomain(chatResponseDto), nil
}

func caseStatus(statusCode int) error {
	//some error validation
	return nil
}
