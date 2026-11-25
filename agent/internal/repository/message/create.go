package message_postgres_repository

import (
	"context"

	"github.com/qury53/GC_learn/agent/internal/domain"
)

func (r *MessageRepository) Create(
	ctx context.Context,
	message domain.Message,
) (domain.Message, error)
