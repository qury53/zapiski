package sessions_postgres_repository

import (
	"github.com/jackc/pgx/v5"
	pgx_pool "github.com/qury53/GC_learn/platform/postgres"
)

// колонки сессии в бд
const sessionColumns = `
	id,
	user_id,
	jti,
	refresh_token_hash,
	user_agent,
	ip_address,
	expires_at,
	revoked_at,
	revoked_reason,
	parent_session_id,
	created_at,
	updated_at
`

// SessionRepository отвечает за сессии пользователя
type SessionRepository struct {
	pool *pgx_pool.Pool
}

// NewSessonRepository конструктор SessionRepository
//
// Returns: *SessionRepository

func NewSessionRepository(
	pool *pgx_pool.Pool,
) *SessionRepository {
	return &SessionRepository{
		pool: pool,
	}
}

// scanSession - вспомогательная ф-ция для сканирования row в модель сессии
func scanSession(row pgx.Row) (SessionModel, error) {
	var m SessionModel
	err := row.Scan(
		&m.ID,
		&m.UserID,
		&m.JTI,
		&m.RefreshTokenHash,
		&m.UserAgent,
		&m.IPAddress,
		&m.ExpiresAt,
		&m.RevokedAt,
		&m.RevokedReason,
		&m.ParentSessionID,
		&m.CreatedAt,
		&m.UpdatedAt,
	)

	return m, err
}
