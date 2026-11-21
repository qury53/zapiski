package users_postgres_repository

import (
	"context"
	"fmt"

	core_errors "github.com/qury53/GC_learn/sso/internal/core/errors"
)

func (r *UserRepository) DeletePendingUser(
	ctx context.Context,
	userID int64,
) error {
	query := `
	DELETE FROM zapiski.users
	WHERE id = $1 AND state = 'pending_verification'
	`

	cmdTag, err := r.pool.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("execute err: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return core_errors.ErrNotFound
	}

	return nil
}
