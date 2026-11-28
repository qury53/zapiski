package user_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/qury53/GC_learn/user/internal/core/domain"
	core_errors "github.com/qury53/GC_learn/user/internal/core/errors"
)

func (r *UserRepository) GetUserByID(
	ctx context.Context,
	userID int64,
) (domain.UserProfile, error) {
	const query = `SELECT ` + userColumns + ` FROM zapiski.users WHERE id = $1`

	row := r.pool.QueryRow(ctx, query, userID)

	var userModel UserProfileModel
	if err := row.Scan(
		&userModel.ID,
		&userModel.Username,
		&userModel.Email,
		&userModel.State,
		&userModel.CreatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.UserProfile{}, fmt.Errorf("GetUserByID repository: %w", core_errors.ErrNotFound)
		}

		return domain.UserProfile{}, fmt.Errorf("GetUserByID repository: scan error; %w", err)
	}

	return modelToDomain(userModel), nil
}
