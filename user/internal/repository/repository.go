package user_repository

import pgx_pool "github.com/qury53/GC_learn/platform/postgres"

type UserRepository struct {
	pool *pgx_pool.Pool
}

func NewUserRepository(pool *pgx_pool.Pool) *UserRepository {
	return &UserRepository{
		pool: pool,
	}
}
