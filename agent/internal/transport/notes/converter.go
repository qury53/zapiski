package notes_transport

import (
	agentv1 "shared/pkg/proto/agent/v1"

	"github.com/qury53/GC_learn/agent/internal/core/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToProtoNote конвертирует доменный Note в protobuf-сообщение.
func ToProtoNote(n domain.Note) *agentv1.Note {
	return &agentv1.Note{
		Id:              n.ID,
		Title:           n.Title,
		Content:         n.Content,
		Source:          toProtoNoteSource(n.Source),
		ChatId:          n.ChatID,
		SourceMessageId: n.SourceMessageID,
		Version:         n.Version,
		CreatedAt:       timestamppb.New(n.CreatedAt),
		UpdatedAt:       timestamppb.New(n.UpdatedAt),
	}
}

// toProtoNoteSource маппит доменный источник конспекта в protobuf enum.
func toProtoNoteSource(source domain.NoteSource) agentv1.NoteSource {
	switch source {
	case domain.NoteSourceManual:
		return agentv1.NoteSource_NOTE_SOURCE_MANUAL
	case domain.NoteSourceMessage:
		return agentv1.NoteSource_NOTE_SOURCE_MESSAGE
	default:
		return agentv1.NoteSource_NOTE_SOURCE_UNSPECIFIED
	}
}

// ToCreateNoteResponse собирает ответ на CreateNote.
func ToCreateNoteResponse(n domain.Note) *agentv1.CreateNoteResponse {
	return &agentv1.CreateNoteResponse{
		Note: ToProtoNote(n),
	}
}

// ToSaveNoteFromMessageResponse собирает ответ на SaveNoteFromMessage.
func ToSaveNoteFromMessageResponse(n domain.Note) *agentv1.SaveNoteFromMessageResponse {
	return &agentv1.SaveNoteFromMessageResponse{
		Note: ToProtoNote(n),
	}
}

// ToGetNoteResponse собирает ответ на GetNote.
func ToGetNoteResponse(n domain.Note) *agentv1.GetNoteResponse {
	return &agentv1.GetNoteResponse{
		Note: ToProtoNote(n),
	}
}

// ToUpdateNoteResponse собирает ответ на UpdateNote.
func ToUpdateNoteResponse(n domain.Note) *agentv1.UpdateNoteResponse {
	return &agentv1.UpdateNoteResponse{
		Note: ToProtoNote(n),
	}
}

// ToListNotesResponse конвертирует слайс доменных конспектов в protobuf-ответ.
func ToListNotesResponse(notes []domain.Note, total int64) *agentv1.ListNotesResponse {
	pbNotes := make([]*agentv1.Note, 0, len(notes))
	for _, note := range notes {
		pbNotes = append(pbNotes, ToProtoNote(note))
	}

	return &agentv1.ListNotesResponse{
		Notes: pbNotes,
		Total: total,
	}
}
