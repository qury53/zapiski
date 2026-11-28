package notes_repository

import (
	"time"

	"github.com/qury53/GC_learn/agent/internal/core/domain"
)

const noteColumns = `
	id,
	version,
	user_id,
	chat_id,
	source_message_id,
	source,
	title,
	content,
	created_at,
	updated_at,
	deleted_at
`

type NoteModel struct {
	ID              int64      `db:"id"`
	UserID          int64      `db:"user_id"`
	Version         int64      `db:"version"`
	Title           string     `db:"title"`
	Content         string     `db:"content"`
	Source          string     `db:"source"`
	ChatID          *int64     `db:"chat_id"`
	SourceMessageID *int64     `db:"source_message_id"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at"`
	DeletedAt       *time.Time `db:"deleted_at"`
}

func modelToDomain(model NoteModel) domain.Note {
	return domain.Note{
		ID:              model.ID,
		UserID:          model.UserID,
		Version:         model.Version,
		Title:           model.Title,
		Content:         model.Content,
		Source:          modelSourceToDomain(model.Source),
		ChatID:          model.ChatID,
		SourceMessageID: model.SourceMessageID,
		CreatedAt:       model.CreatedAt,
		UpdatedAt:       model.UpdatedAt,
		DeletedAt:       model.DeletedAt,
	}
}

func modelSourceToDomain(source string) domain.NoteSource {
	switch source {
	case "manual":
		return domain.NoteSourceManual
	case "message":
		return domain.NoteSourceMessage
	default:
		return domain.NoteSourceUnspecified
	}
}
