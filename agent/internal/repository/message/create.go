package message_postgres_repository

import (
	"context"
	"fmt"

	"github.com/qury53/GC_learn/agent/internal/core/domain"
)

func (r *MessageRepository) Create(
	ctx context.Context,
	message domain.Message,
) (domain.Message, error) {
	const query = `
	INSERT INTO zapiski.messages (
		chat_id,
		role,
		content
	)
	VALUES ($1, $2, $3)
	RETURNING ` + messageColumns

	row := r.pool.QueryRow(ctx, query, message.ChatID, message.Role, message.Content)
	var messageModel MessageModel
	if err := row.Scan(
		&messageModel.ID,
		&messageModel.ChatID,
		&messageModel.Role,
		&messageModel.Content,
		&messageModel.CreatedAt,
	); err != nil {
		return domain.Message{}, fmt.Errorf("insert message error: %w", err)
	}

	return modelToDomain(messageModel), nil
}
