package sessions_postgres_repository

import (
	"context"
	"fmt"

	"github.com/qury53/GC_learn/sso/internal/core/domain"
)

func (r *SessionRepository) ListActiveByUserID(ctx context.Context, userID int64) ([]domain.Session, error) {
	const query = `
		SELECT ` + sessionColumns + `
		FROM zapiski.user_sessions
		WHERE user_id = $1
		  AND revoked_at IS NULL
		  AND expires_at > now()
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, userID)

	if err != nil {
		return nil, fmt.Errorf("session_repository.ListActiveByUserID: %w", err)
	}

	defer rows.Close()

	var result []domain.Session

	for rows.Next() {
		m, err := scanSession(rows)
		if err != nil {
			return nil, fmt.Errorf("session_repository.ListActiveByUserID scan: %w", err)
		}

		result = append(result, modelToDomain(m))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("session_repository.ListActiveByUserID rows: %w", err)
	}

	return result, nil
}
