package notes_transport

import (
	"context"
	agentv1 "shared/pkg/proto/agent/v1"

	"github.com/qury53/GC_learn/agent/internal/core/domain"
	"github.com/qury53/GC_learn/platform/logger"
)

type NoteTransport struct {
	agentv1.UnimplementedNoteServiceServer
	log         *logger.Logger
	noteService NoteService
}

type NoteService interface {
	CreateNote(ctx context.Context, note domain.Note) (domain.Note, error)
	SaveNoteFromMessage(ctx context.Context, note domain.Note) (domain.Note, error)
}
