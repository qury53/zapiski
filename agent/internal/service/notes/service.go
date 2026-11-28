package notes_service

import (
	"context"

	"github.com/qury53/GC_learn/agent/internal/core/domain"
	"github.com/qury53/GC_learn/platform/logger"
)

type NoteService struct {
	log            *logger.Logger
	noteRepository NoteRepository
}

type ChatRepository interface {
}

type NoteRepository interface {
	CreateNote(ctx context.Context, note domain.Note) (domain.Note, error)
	SaveNoteFromMessage(ctx context.Context, note domain.Note) (domain.Note, error)
}

func NewNoteService(
	log *logger.Logger,
	noteRepository NoteRepository,
) *NoteService {
	return &NoteService{
		log:            log,
		noteRepository: noteRepository,
	}
}
