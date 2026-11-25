package message_postgres_repository

import pgx_pool "github.com/qury53/GC_learn/platform/postgres"

type MessageRepository struct {
	pool *pgx_pool.Pool
}

func NewMessageRepository(pool *pgx_pool.Pool) *MessageRepository {
	return &MessageRepository{
		pool: pool,
	}
}
