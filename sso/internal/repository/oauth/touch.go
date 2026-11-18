package oauth_postgres_repository

import (
	"context"
	"fmt"
)

// Touch обновляет email и флаг подтверждения у привязки,
// если провайдер вернул изменённые данные.
//
// Returns: error

func (r *OAuthRepository) Touch(
	ctx context.Context,
	id int64,
	email string,
	emailVerified bool,
) error {
	const query = `
		UPDATE zapiski.oauth_accounts
		SET email = $2, email_verified = $3, updated_at = now()
		WHERE id = $1
	`

	if _, err := r.pool.Exec(ctx, query, id, nullableString(email), emailVerified); err != nil {
		return fmt.Errorf("oauth_repository.Touch: %w", err)
	}

	return nil
}
