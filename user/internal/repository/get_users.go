package user_repository

import (
	"context"
	"fmt"

	"github.com/qury53/GC_learn/user/internal/core/domain"
)

func (r *UserRepository) GetUsers(
	ctx context.Context,
	limit *int64,
	offset *int64,
) ([]domain.UserProfile, error) {
	const query = `SELECT ` + userColumns + ` FROM zapiski.users
	ORDER BY id DESC
	LIMIT $1
	OFFSET $2 
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("GetUsers repository: query err: %w", err)
	}

	defer rows.Close()

	var userModels []UserProfileModel
	for rows.Next() {
		var userModel UserProfileModel
		if err := rows.Scan(
			&userModel.ID,
			&userModel.Username,
			&userModel.Email,
			&userModel.State,
			&userModel.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf(
				"GetUsers repository: scan error: %w",
				err,
			)
		}

		userModels = append(userModels, userModel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetUsers repository: rows error: %w", err)
	}

	return modelsToDomain(userModels), nil
}
