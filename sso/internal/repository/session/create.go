package sessions_postgres_repository

import (
	"context"
	"fmt"

	"github.com/qury53/GC_learn/sso/internal/core/domain"
)

func (r *SessionRepository) Create(
	ctx context.Context,
	s domain.Session,
) (domain.Session, error) {
	const query = `
	INSERT INTO zapiski.user_sessions (
		user_id,
		jti,
		refresh_token_hash,
		user_agent,
		ip_address,
		expires_at,
		parent_session_id
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING ` + sessionColumns

	row := r.pool.QueryRow(
		ctx,
		query,
		s.UserID,
		s.JTI,
		s.RefreshTokenHash,
		s.UserAgent,
		s.IPAddress,
		s.ExpiresAt,
		s.ParentSessionID,
	)

	m, err := scanSession(row)
	if err != nil {
		return domain.Session{}, fmt.Errorf("sessions.Create scan error: %w", err)
	}

	return modelToDomain(m), nil
}
