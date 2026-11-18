package sessions_postgres_repository

import (
	"context"
	"fmt"

	core_errors "github.com/qury53/GC_learn/sso/internal/core/errors"
)

func (r *SessionRepository) Revoke(ctx context.Context, jti string, reason string) error {
	const query = `
		UPDATE zapiski.user_sessions
		SET revoked_at = now(), revoked_reason = $2
		WHERE jti = $1 AND revoked_at IS NULL
	`
	tag, err := r.pool.Exec(ctx, query, jti, reason)
	if err != nil {
		return fmt.Errorf("session_repository.Revoke: %w", err)
	}

	if tag.RowsAffected() == 0 {
		var exists bool
		err := r.pool.QueryRow(ctx,
			`SELECT EXISTS(SELECT 1 FROM zapiski.user_sessions WHERE jti = $1)`, jti,
		).Scan(&exists)

		if err != nil {
			return fmt.Errorf("session_repository.Revoke exists: %w", err)
		}

		if !exists {
			return core_errors.ErrNotFound
		}

		return core_errors.ErrSessionRevoked
	}

	return nil
}
