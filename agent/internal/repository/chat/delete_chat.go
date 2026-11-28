package chat_postgres_repository

import (
	"context"
	"fmt"

	core_errors "github.com/qury53/GC_learn/agent/internal/core/errors"
)

func (r *ChatRepository) DeleteChat(
	ctx context.Context,
	chatID int,
) error {
	query := `
		DELETE FROM zapiski.chats
		WHERE id = $1;
	`

	cmdTag, err := r.pool.Exec(ctx, query, chatID)
	if err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("delete: chat not found: %w", core_errors.ErrNotFound)
	}

	return nil
}
