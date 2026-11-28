package notes_repository

import (
	"context"
	"fmt"

	"github.com/qury53/GC_learn/agent/internal/core/domain"
)

func (r *NoteRepository) CreateNote(
	ctx context.Context,
	note domain.Note,
) (domain.Note, error) {
	const query = `INSERT INTO zapiski.notes (
		title,
		content,
		chat_id,
		user_id
	)
	VALUES ($1, $2, $3, $4)
	RETURNING ` + noteColumns

	row := r.pool.QueryRow(
		ctx,
		query,
		note.Title,
		note.Content,
		note.ChatID,
		note.UserID,
	)

	var noteModel NoteModel

	if err := row.Scan(
		&noteModel.ID,
		&noteModel.Version,
		&noteModel.UserID,
		&noteModel.ChatID,
		&noteModel.SourceMessageID,
		&noteModel.Source,
		&noteModel.Title,
		&noteModel.Content,
		&noteModel.CreatedAt,
		&noteModel.UpdatedAt,
		&noteModel.DeletedAt,
	); err != nil {
		return domain.Note{}, fmt.Errorf("CreateNote repository: scan error: %w", err)
	}

	return modelToDomain(noteModel), nil
}
