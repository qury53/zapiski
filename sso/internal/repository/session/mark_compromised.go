package sessions_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	core_errors "github.com/qury53/GC_learn/sso/internal/core/errors"
)

func (r *SessionRepository) MarkCompromised(ctx context.Context, jti string, reason string) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("session_repository.MarkCompromised begin: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	const findRoot = `
		WITH RECURSIVE chain AS (
			SELECT id, parent_session_id, user_id
			FROM zapiski.user_sessions
			WHERE jti = $1
			UNION ALL
			SELECT s.id, s.parent_session_id, s.user_id
			FROM zapiski.user_sessions s
			JOIN chain c ON s.id = c.parent_session_id
		)
		SELECT id, user_id
		FROM chain
		WHERE parent_session_id IS NULL
		LIMIT 1
	`
	var rootID, userID int64
	if err := tx.QueryRow(ctx, findRoot, jti).Scan(&rootID, &userID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core_errors.ErrNotFound
		}
		return fmt.Errorf("session_repository.MarkCompromised root: %w", err)
	}

	const revokeChain = `
		WITH RECURSIVE chain AS (
			SELECT id
			FROM zapiski.user_sessions
			WHERE id = $1
			UNION ALL
			SELECT s.id
			FROM zapiski.user_sessions s
			JOIN chain c ON s.parent_session_id = c.id
		)
		UPDATE zapiski.user_sessions
		SET revoked_at = now(), revoked_reason = $2
		WHERE id IN (SELECT id FROM chain) AND revoked_at IS NULL
	`
	if _, err := tx.Exec(ctx, revokeChain, rootID, reason); err != nil {
		return fmt.Errorf("session_repository.MarkCompromised revoke: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("session_repository.MarkCompromised commit: %w", err)
	}
	return nil
}
