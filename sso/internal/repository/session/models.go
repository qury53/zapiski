package sessions_postgres_repository

import (
	"net/netip"
	"time"

	"github.com/qury53/GC_learn/sso/internal/core/domain"
)

type SessionModel struct {
	ID               int64
	UserID           int64
	JTI              string
	RefreshTokenHash string
	UserAgent        *string
	IPAddress        *netip.Addr
	ExpiresAt        time.Time
	RevokedAt        *time.Time
	RevokedReason    *string
	ParentSessionID  *int64
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func modelToDomain(model SessionModel) domain.Session {
	return domain.Session{
		ID:               model.ID,
		UserID:           model.UserID,
		JTI:              model.JTI,
		RefreshTokenHash: model.RefreshTokenHash,
		UserAgent:        model.UserAgent,
		IPAddress:        model.IPAddress,
		ExpiresAt:        model.ExpiresAt,
		RevokedAt:        model.RevokedAt,
		RevokedReason:    model.RevokedReason,
		ParentSessionID:  model.ParentSessionID,
		CreatedAt:        model.CreatedAt,
		UpdatedAt:        model.UpdatedAt,
	}
}
