package chat_postgres_repository

import (
	"context"
	"fmt"

	"github.com/qury53/GC_learn/agent/internal/domain"
)

func (r *ChatRepository) GetChats(
	ctx context.Context,
	userID int,
	limit int,
	offset int,
) ([]domain.Chat, error) {
	const query = `
	SELECT ` + chatColumns + `
	FROM zapiski.chats
	WHERE user_id = $1 AND deleted_at IS NULL
	ORDER BY updated_at DESC
	LIMIT $2
	OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	defer rows.Close()

	chatModels := make([]ChatModel, 0, limit)
	for rows.Next() {
		var chatModel ChatModel
		if err := rows.Scan(
			&chatModel.ID,
			&chatModel.Version,
			&chatModel.UserID,
			&chatModel.Mode,
			&chatModel.CreatedAt,
			&chatModel.UpdatedAt,
			&chatModel.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		chatModels = append(chatModels, chatModel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return modelsToDomain(chatModels)
}
