package users_postgres_repository

import (
	"context"
	"fmt"

	core_errors "github.com/qury53/GC_learn/sso/internal/core/errors"
)

func (r *UserRepository) ActivateUser(
	ctx context.Context,
	email string,
) error {
	query := `
	UPDATE zapiski.users
	SET state = 'active',
		updated_at = now(),
		version = version+1
	WHERE email = $1;
	`

	cmdTag, err := r.pool.Exec(ctx, query, email)
	if err != nil {
		return fmt.Errorf("user activation failed: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("user activation failed: %w", core_errors.ErrNotFound)
	}

	return nil
}
