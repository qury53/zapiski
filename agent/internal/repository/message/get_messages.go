package message_postgres_repository

import (
	"context"
	"fmt"

	"github.com/qury53/GC_learn/agent/internal/core/domain"
)

func (r *MessageRepository) GetMessages(
	ctx context.Context,
	chatID int,
	limit int,
	offset int,
) ([]domain.Message, error) {
	const query = `
		SELECT ` + messageColumns + `
		FROM (
			SELECT ` + messageColumns + `
			FROM zapiski.messages
			WHERE chat_id = $1
			ORDER BY created_at DESC, id DESC
			LIMIT $2
			OFFSET $3
		) AS recent
		ORDER BY created_at ASC, id ASC
	`
	rows, err := r.pool.Query(ctx, query, chatID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	defer rows.Close()

	messages := make([]MessageModel, 0, limit)
	for rows.Next() {
		var messageModel MessageModel
		if err := rows.Scan(
			&messageModel.ID,
			&messageModel.ChatID,
			&messageModel.Role,
			&messageModel.Content,
			&messageModel.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		messages = append(messages, messageModel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return modelsToDomain(messages), nil
}
