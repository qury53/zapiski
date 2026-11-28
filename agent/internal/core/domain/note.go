package domain

import "time"

type NoteSource int32

const (
	NoteSourceUnspecified NoteSource = iota // 0
	NoteSourceManual                        // 1
	NoteSourceMessage                       // 2
)

type Note struct {
	ID              int64
	UserID          int64
	Version         int64
	Title           string
	Content         string
	Source          NoteSource
	ChatID          *int64
	SourceMessageID *int64
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
}

// Конструктор для создания uninitialized конспекта, нужен будет для
// ручного сохранения конспекта
func NewNoteUninitialized(
	title string,
	content string,
	chatID *int64,
	userID int64,
) Note {
	now := time.Now()

	return Note{
		ID:        NoteIDUninitialized,
		UserID:    userID,
		ChatID:    chatID,
		Version:   VersionUnitialized,
		Title:     title,
		Content:   content,
		Source:    NoteSourceManual,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Конструкитор для создания uninitialized конспекта из сообщения агента
func NewNoteFromMessage(userID int64, msg Message, title string) Note {
	now := time.Now()

	chatID := int64(msg.ChatID)
	messageID := int64(msg.ID)

	return Note{
		ID:              NoteIDUninitialized,
		UserID:          userID,
		Title:           title,
		Content:         msg.Content,
		Source:          NoteSourceMessage,
		ChatID:          &chatID,
		SourceMessageID: &messageID,
		Version:         VersionUnitialized,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}
