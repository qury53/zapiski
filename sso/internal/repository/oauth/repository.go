package oauth_postgres_repository

import pgx_pool "github.com/qury53/GC_learn/platform/postgres"

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
	pool *pgx_pool.Pool
}

// NewOAuthRepository конструктор OAuthRepository
//
// Returns: *OAuthRepository

func NewOAuthRepository(
	pool *pgx_pool.Pool,
) *OAuthRepository {
	return &OAuthRepository{
		pool: pool,
	}
}
