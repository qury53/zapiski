package ai_provider_gc_models

import (
	"time"

	"github.com/qury53/GC_learn/agent/internal/domain"
)

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at"`
}

func AuthModelToDomain(model AuthResponse) domain.AuthResponse {
	return domain.AuthResponse{
		AccessToken: model.AccessToken,
		ExpiredAt:   time.UnixMilli(model.ExpiresAt),
	}
}
