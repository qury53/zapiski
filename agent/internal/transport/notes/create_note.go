package notes_transport

import (
	"context"

	agentv1 "shared/pkg/proto/agent/v1"

	"github.com/qury53/GC_learn/agent/internal/core/domain"
	core_errors "github.com/qury53/GC_learn/agent/internal/core/errors"
	"github.com/qury53/GC_learn/platform/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (t *NoteTransport) CreateNote(
	ctx context.Context,
	req *agentv1.CreateNoteRequest,
) (*agentv1.CreateNoteResponse, error) {
	userID, ok := auth.UserIDFrom(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "no identity")
	}

	note := domain.NewNoteUninitialized(
		req.Title,
		req.Content,
		req.ChatId,
		userID,
	)

	note, err := t.noteService.CreateNote(ctx, note)
	if err != nil {
		return nil, status.Error(
			core_errors.RespondCode(err),
			"CreateNote transport: failed to create note",
		)
	}

	return ToCreateNoteResponse(note), nil
}
