package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/qury53/GC_learn/sso/internal/core/domain"
	core_errors "github.com/qury53/GC_learn/sso/internal/core/errors"
)

func (r *UserRepository) SaveUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	query := `
		INSERT INTO zapiski.users (username, password, email)
		VALUES ($1, $2, $3)
		RETURNING id, version, username, password, email, state, created_at, updated_at;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		user.Username,
		user.Password,
		user.Email,
	)

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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.User{}, fmt.Errorf("user already exists: %w", core_errors.ErrConflict)
		}

		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	return modelToDomain(userModel), nil
}
