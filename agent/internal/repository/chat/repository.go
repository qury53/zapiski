package chat_postgres_repository

import pgx_pool "github.com/qury53/GC_learn/platform/postgres"

type ChatRepository struct {
	pool *pgx_pool.Pool
}

func NewChatRepository(pool *pgx_pool.Pool) *ChatRepository {
	return &ChatRepository{
		pool: pool,
	}
}
