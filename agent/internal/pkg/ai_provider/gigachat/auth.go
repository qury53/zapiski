package ai_provider_gc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/qury53/GC_learn/agent/internal/domain"
	errors "github.com/qury53/GC_learn/agent/internal/errors"
	ai_provider_gc_models "github.com/qury53/GC_learn/agent/pkg/ai_provider/gigachat/models"
)

func (p *GCProvider) requestAccessToken(ctx context.Context) (domain.AuthResponse, error) {
	url := "https://ngw.devices.sberbank.ru:9443/api/v2/oauth"
	payload := strings.NewReader(fmt.Sprintf("scope=%s", p.scope))

	req, err := http.NewRequestWithContext(ctx, "POST", url, payload)
	if err != nil {
		return domain.AuthResponse{}, fmt.Errorf("auth error: %w", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("RqUID", uuid.NewString())
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", p.apiKey))
	resp, err := p.client.Do(req)

	if err != nil {
		return domain.AuthResponse{}, fmt.Errorf("http client error: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return domain.AuthResponse{}, fmt.Errorf("failed to read response body: %w", err)
	}

	switch resp.StatusCode {
	case http.StatusBadRequest:
		return domain.AuthResponse{}, fmt.Errorf("%w: %s", errors.ErrInvalidArgument, string(body))
	case http.StatusUnauthorized:
		return domain.AuthResponse{}, fmt.Errorf("%w: %s", errors.ErrUnauthorized, string(body))
	}

	var response ai_provider_gc_models.AuthResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return domain.AuthResponse{}, fmt.Errorf("failed to decode: %w", errors.ErrInvalidArgument)
	}

	p.accessToken = response.AccessToken
	p.expiresAt = time.UnixMilli(response.ExpiresAt)
	return ai_provider_gc_models.AuthModelToDomain(response), nil
}

func (p *GCProvider) refreshAccessToken(ctx context.Context) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.isValid(p.accessToken, p.expiresAt) {
		return p.accessToken, nil
	}

	authResp, err := p.requestAccessToken(ctx)
	if err != nil {
		return "", fmt.Errorf("refresh gigachat access token: %w", err)
	}

	if authResp.AccessToken == "" || authResp.ExpiredAt.Before(time.Now()) {
		return "", fmt.Errorf("gigachat auth response has empty token or expiry")
	}

	p.accessToken = authResp.AccessToken
	p.expiresAt = authResp.ExpiredAt
	return p.accessToken, nil
}

func (p *GCProvider) isValid(token string, expiresAt time.Time) bool {
	return token != "" && time.Now().Before(expiresAt.Add(-30*time.Second))
}
