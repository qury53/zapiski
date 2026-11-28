package message_postgres_repository

import (
	"time"

	"github.com/qury53/GC_learn/agent/internal/core/domain"
)

const (
	messageColumns = `
	id,
	chat_id,
	role,
	content,
	created_at
	`
)

type MessageModel struct {
	CreatedAt time.Time
	Role      string
	Content   string
	ID        int
	ChatID    int
}

func modelToDomain(model MessageModel) domain.Message {
	return domain.Message{
		ID:        model.ID,
		ChatID:    model.ChatID,
		Role:      domain.Role(model.Role),
		Content:   model.Content,
		CreatedAt: model.CreatedAt,
	}
}

func modelsToDomain(models []MessageModel) []domain.Message {
	domains := make([]domain.Message, len(models))
	for i, model := range models {
		domains[i] = modelToDomain(model)
	}

	return domains
}
