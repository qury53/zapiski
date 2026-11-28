package chat_postgres_repository

import (
	"context"
	"fmt"

	"github.com/qury53/GC_learn/agent/internal/core/domain"
)

func (r *ChatRepository) CreateChat(
	ctx context.Context,
	chat domain.Chat,
) (domain.Chat, error) {
	const query = `
	INSERT INTO zapiski.chats (
		user_id, 
		mode
	)
	VAlUES ($1, $2)
	RETURNING ` + chatColumns

	mode, err := mapDomainModeToModel(chat.Mode)
	if err != nil {
		return domain.Chat{}, err
	}

	row := r.pool.QueryRow(ctx, query, chat.UserID, mode)

	var chatModel ChatModel

	if err := row.Scan(
		&chatModel.ID,
		&chatModel.Version,
		&chatModel.UserID,
		&chatModel.Mode,
		&chatModel.CreatedAt,
		&chatModel.UpdatedAt,
		&chatModel.DeletedAt,
	); err != nil {
		return domain.Chat{}, fmt.Errorf("scan error: %w", err)
	}

	return modelToDomain(chatModel)
}
