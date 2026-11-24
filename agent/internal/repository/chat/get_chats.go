package chat_postgres_repository

import (
	"context"
	"fmt"

	"github.com/qury53/GC_learn/agent/internal/domain"
)

func (r *ChatRepository) GetChats(
	ctx context.Context,
	limit *int,
	offset *int,
) ([]domain.Chat, error) {
	const query = `
	SELECT ` + chatColumns + ` FROM zapiski.chats
	LIMIT $1
	OFFSET $2
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	var chatModels []ChatModel
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
			return nil, fmt.Errorf("scan error")
		}

		chatModels = append(chatModels, chatModel)
	}

	return modelsToDomain(chatModels)
}
