package notes_repository

import pgx_pool "github.com/qury53/GC_learn/platform/postgres"

type NoteRepository struct {
	pool *pgx_pool.Pool
}

func NewNoteRepository(pool *pgx_pool.Pool) *NoteRepository {
	return &NoteRepository{
		pool: pool,
	}
}
