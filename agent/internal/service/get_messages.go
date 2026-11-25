package agent_service

import (
	"context"
	"fmt"

	"github.com/qury53/GC_learn/agent/internal/domain"
	core_errors "github.com/qury53/GC_learn/agent/internal/errors"
)

const (
	defaultMessagesLimit = 100
	maxMessagesLimit     = 500
)

func (s *AgentService) GetMessages(
	ctx context.Context,
	chatID int,
	userID int,
	limit *int,
	offset *int,
) ([]domain.Message, error) {
	if userID <= 0 {
		return nil, fmt.Errorf(
			"GetMessages service: userID has to be positive: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	l := defaultMessagesLimit
	if limit != nil {
		l = *limit
	}

	if l <= 0 {
		return nil, fmt.Errorf(
			"GetMessages service: limit has to be positive: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if l > maxMessagesLimit {
		l = maxMessagesLimit
	}

	o := 0
	if offset != nil {
		o = *offset
	}

	if o < 0 {
		return nil, fmt.Errorf(
			"GetMessages service: offset cannot be negative: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	// Чужой чат отдаём как отсутствующий: не палим его существование.
	chat, err := s.chatRepository.GetChat(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("GetMessages service: failed to get chat: %w", err)
	}

	if chat.UserID != userID {
		return nil, fmt.Errorf(
			"GetMessages service: chat not found: %w",
			core_errors.ErrNotFound,
		)
	}

	messages, err := s.messageRepository.GetMessages(ctx, chatID, l, o)
	if err != nil {
		return nil, fmt.Errorf("GetMessages service: failed to get messages: %w", err)
	}

	return messages, nil
}
