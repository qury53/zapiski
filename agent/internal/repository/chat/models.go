package chat_postgres_repository

import (
	"fmt"
	"time"

	"github.com/qury53/GC_learn/agent/internal/core/domain"
)

const chatColumns = `
	id,
	version,
	user_id,
	mode,
	created_at,
	updated_at,
	deleted_at
`

type ChatModel struct {
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
	Mode      string     `db:"mode"`
	ID        int        `db:"id"`
	Version   int        `db:"version"`
	UserID    int        `db:"user_id"`
}

func modelToDomain(model ChatModel) (domain.Chat, error) {
	domainMode, err := mapModelModeToDomain(model.Mode)
	if err != nil {
		return domain.Chat{}, err
	}

	return domain.Chat{
		ID:        model.ID,
		Version:   model.Version,
		Mode:      domainMode,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		DeletedAt: model.DeletedAt,
		UserID:    model.UserID,
	}, nil
}

func modelsToDomain(models []ChatModel) ([]domain.Chat, error) {
	domainChats := make([]domain.Chat, len(models))
	for i, model := range models {
		chat, err := modelToDomain(model)
		if err != nil {
			return nil, fmt.Errorf("convert chat model %d (index %d): %w", model.ID, i, err)
		}

		domainChats[i] = chat
	}

	return domainChats, nil
}

func mapModelModeToDomain(mode string) (domain.ChatMode, error) {
	switch mode {
	case "full":
		return domain.ModeFull, nil
	case "support":
		return domain.ModeSupport, nil
	default:
		// Неизвестный режим не должен ронять всю выдачу списка чатов.
		return domain.ModeUnknown, nil
	}
}

func mapDomainModeToModel(mode domain.ChatMode) (string, error) {
	switch mode {
	case domain.ModeSupport:
		return "support", nil
	case domain.ModeFull, domain.ModeUnknown:
		return "full", nil
	default:
		return "full", nil
	}
}
