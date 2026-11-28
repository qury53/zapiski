package chat_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/qury53/GC_learn/agent/internal/core/domain"
	core_errors "github.com/qury53/GC_learn/agent/internal/core/errors"
)

// GetChat возвращает чат по id. Используется для проверки владельца перед
// выдачей сообщений и записью в историю.
func (r *ChatRepository) GetChat(
	ctx context.Context,
	chatID int,
) (domain.Chat, error) {
	const query = `
	SELECT ` + chatColumns + `
	FROM zapiski.chats
	WHERE id = $1 AND deleted_at IS NULL
	`

	var chatModel ChatModel

	err := r.pool.QueryRow(ctx, query, chatID).Scan(
		&chatModel.ID,
		&chatModel.Version,
		&chatModel.UserID,
		&chatModel.Mode,
		&chatModel.CreatedAt,
		&chatModel.UpdatedAt,
		&chatModel.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Chat{}, fmt.Errorf("chat not found: %w", core_errors.ErrNotFound)
		}

		return domain.Chat{}, fmt.Errorf("scan error: %w", err)
	}

	return modelToDomain(chatModel)
}
