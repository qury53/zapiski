package notes_service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/qury53/GC_learn/agent/internal/core/domain"
)

func (s *NoteService) CreateNote(
	ctx context.Context,
	note domain.Note,
) (domain.Note, error) {
	if note.ChatID != nil {

	}
	note, err := s.noteRepository.CreateNote(ctx, note)
	if err != nil {
		s.log.Debug(
			"CreateNote service: failed to create note",
			slog.Any("error", err),
		)

		return domain.Note{}, fmt.Errorf("CreateNote service: %w", err)
	}

	return note, nil
}
