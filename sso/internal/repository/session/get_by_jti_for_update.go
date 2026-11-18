package sessions_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/qury53/GC_learn/sso/internal/core/domain"
	core_errors "github.com/qury53/GC_learn/sso/internal/core/errors"
)

func (r *SessionRepository) GetByJTIForUpdate(ctx context.Context, jti string) (domain.Session, error) {
	const query = `SELECT ` + sessionColumns + ` FROM zapiski.user_sessions WHERE jti = $1 FOR UPDATE`

	m, err := scanSession(r.pool.QueryRow(ctx, query, jti))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Session{}, core_errors.ErrNotFound
		}

		return domain.Session{}, fmt.Errorf("session_repository.GetByJTIForUpdate: %w", err)
	}

	return modelToDomain(m), nil
}
