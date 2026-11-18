package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/qury53/GC_learn/sso/internal/core/domain"
	core_errors "github.com/qury53/GC_learn/sso/internal/core/errors"
)

// GetUserByID возвращает пользователя по первичному ключу.
//
// Returns: domain.User, error
func (r *UserRepository) GetUserByID(
	ctx context.Context,
	id int64,
) (domain.User, error) {
	const query = `
		SELECT id, version, username, password, email, state, created_at, updated_at
		FROM zapiski.users
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)

	var userModel UserModel
	if err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.Username,
		&userModel.Password,
		&userModel.Email,
		&userModel.State,
		&userModel.CreatedAt,
		&userModel.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf("user not found: %w", core_errors.ErrNotFound)
		}

		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	return modelToDomain(userModel), nil
}
