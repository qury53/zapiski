package sessions_postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/qury53/GC_learn/sso/internal/core/domain"
	core_errors "github.com/qury53/GC_learn/sso/internal/core/errors"
)

func (r *SessionRepository) Rotate(
	ctx context.Context,
	jti string,
	params domain.RotateParams,
) (domain.Session, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return domain.Session{}, fmt.Errorf("session_repository.Rotate begin: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	const selectQuery = `
		SELECT ` + sessionColumns + `
		FROM zapiski.user_sessions
		WHERE jti = $1
		FOR UPDATE
	`
	old, err := scanSession(tx.QueryRow(ctx, selectQuery, jti))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Session{}, core_errors.ErrNotFound
		}
		return domain.Session{}, fmt.Errorf("session_repository.Rotate select: %w", err)
	}

	if old.RevokedAt != nil {
		return domain.Session{}, core_errors.ErrSessionRevoked
	}
	if !old.ExpiresAt.After(time.Now()) {
		return domain.Session{}, core_errors.ErrSessionExpired
	}

	const revokeQuery = `
		UPDATE zapiski.user_sessions
		SET revoked_at = now(), revoked_reason = 'rotated'
		WHERE id = $1
	`
	if _, err := tx.Exec(ctx, revokeQuery, old.ID); err != nil {
		return domain.Session{}, fmt.Errorf("session_repository.Rotate revoke: %w", err)
	}

	const insertQuery = `
		INSERT INTO zapiski.user_sessions (
			user_id, jti, refresh_token_hash, user_agent, ip_address,
			expires_at, parent_session_id
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING ` + sessionColumns

	newRow := tx.QueryRow(ctx, insertQuery,
		old.UserID,
		params.NewJTI,
		params.NewRefreshTokenHash,
		params.UserAgent,
		params.IPAddress,
		params.NewExpiresAt,
		old.ID,
	)

	created, err := scanSession(newRow)
	if err != nil {
		return domain.Session{}, fmt.Errorf("session_repository.Rotate insert: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.Session{}, fmt.Errorf("session_repository.Rotate commit: %w", err)
	}

	return modelToDomain(created), nil
}
