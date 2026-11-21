package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/qury53/GC_learn/sso/internal/core/domain"
	core_errors "github.com/qury53/GC_learn/sso/internal/core/errors"
)

// GetUserByEmail возвращает пользователя по email. Сравнение точное:
// email в БД хранится нормализованным (lower+trim) при регистрации и в OAuth.
//
// Returns: domain.User, error

func (r *UserRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (domain.User, error) {
	const query = `
		SELECT id, version, username, password, email, state, created_at, updated_at
		FROM zapiski.users
		WHERE email = $1
	`

	return r.getUser(ctx, query, email)
}

// GetUserByUsername возвращает найденого пользователя по имени, если его нет,
// то возвращает ошибку ErrNotFound
//
// Returns: domain.User, error

func (r *UserRepository) GetUserByUsername(
	ctx context.Context,
	username string,
) (domain.User, error) {
	const query = `
		SELECT id, version, username, password, email, state, created_at, updated_at
		FROM zapiski.users
		WHERE username = $1;
	`

	return r.getUser(ctx, query, username)
}

// getUser выполняет запрос и мапит строку в доменную сущность.
//
// Returns: domain.User, error

func (r *UserRepository) getUser(
	ctx context.Context,
	query string,
	value string,
) (domain.User, error) {
	row := r.pool.QueryRow(ctx, query, value)
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
