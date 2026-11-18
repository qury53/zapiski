package oauth_postgres_repository

import (
	core_repository_postgres "github.com/qury53/GC_learn/sso/internal/core/repository/postgres"
)

// колонки oauth аккаунта в бд
const oauthAccountColumns = `
	id,
	user_id,
	provider,
	provider_user_id,
	email,
	email_verified,
	created_at,
	updated_at
`

// OAuthRepository отвечает за привязки внешних OAuth аккаунтов
type OAuthRepository struct {
	pool *core_repository_postgres.Pool
}

// NewOAuthRepository конструктор OAuthRepository
//
// Returns: *OAuthRepository

func NewOAuthRepository(
	pool *core_repository_postgres.Pool,
) *OAuthRepository {
	return &OAuthRepository{
		pool: pool,
	}
}
