package domain

import "time"

type ChatMode int

const (
	ModeUnknown ChatMode = iota // 0
	ModeSupport                 // 1
	ModeFull                    // 2
)

type Chat struct {
	ID        int
	Version   int
	Mode      ChatMode
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	UserID    int
}

func NewChatUninitialized(
	userID int,
	mode ChatMode,
) Chat {
	t := time.Now()
	return Chat{
		ID:        ChatIDUninitialized,
		Mode:      mode,
		CreatedAt: t,
		UpdatedAt: t,
		UserID:    userID,
	}
}
