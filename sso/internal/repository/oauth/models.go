package oauth_postgres_repository

import (
	"time"

	"github.com/qury53/GC_learn/sso/internal/core/domain"
)

// OAuthAccountModel модель oauth аккаунта для базы данных
type OAuthAccountModel struct {
	ID             int64
	UserID         int64
	Provider       string
	ProviderUserID string
	Email          *string
	EmailVerified  bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// modelToDomain конвертер модели из базы в доменную сущность
//
// Returns: domain.OAuthAccount

func modelToDomain(m OAuthAccountModel) domain.OAuthAccount {
	var email string
	if m.Email != nil {
		email = *m.Email
	}

	return domain.OAuthAccount{
		ID:             m.ID,
		UserID:         m.UserID,
		Provider:       domain.Provider(m.Provider),
		ProviderUserID: m.ProviderUserID,
		Email:          email,
		EmailVerified:  m.EmailVerified,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}
