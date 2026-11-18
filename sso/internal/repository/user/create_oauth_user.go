package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/qury53/GC_learn/sso/internal/core/domain"
	core_errors "github.com/qury53/GC_learn/sso/internal/core/errors"
)

// CreateOAuthUser создает пользователя с подтвержденным email и без пароля.
// Пароль остается NULL — OAuth-аккаунт не должен авторизовываться через Login.
//
// Returns: domain.User, error
func (r *UserRepository) CreateOAuthUser(
	ctx context.Context,
	username string,
	email string,
) (domain.User, error) {
	const query = `
		INSERT INTO zapiski.users (username, password, email, state)
		VALUES ($1, NULL, $2, 'active')
		RETURNING id, version, username, password, email, state, created_at, updated_at
	`

	var model UserModel
	if err := r.pool.QueryRow(ctx, query, username, email).Scan(
		&model.ID,
		&model.Version,
		&model.Username,
		&model.Password,
		&model.Email,
		&model.State,
		&model.CreatedAt,
		&model.UpdatedAt,
	); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.User{}, fmt.Errorf("user already exists: %w", core_errors.ErrConflict)
		}

		return domain.User{}, fmt.Errorf("user repository.CreateOAuthUser: %w", err)
	}

	return modelToDomain(model), nil
}
