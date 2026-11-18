package sessions_postgres_repository

import (
	"context"
	"fmt"
)

func (r *SessionRepository) RevokeAllForUser(ctx context.Context, userID int64, reason string) error {
	const query = `
		UPDATE zapiski.user_sessions
		SET revoked_at = now(), revoked_reason = $2
		WHERE user_id = $1 AND revoked_at IS NULL
	`
	if _, err := r.pool.Exec(ctx, query, userID, reason); err != nil {
		return fmt.Errorf("session_repository.RevokeAllForUser: %w", err)
	}

	return nil
}
