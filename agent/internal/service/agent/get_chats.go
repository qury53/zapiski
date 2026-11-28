package agent_service

import (
	"context"
	"fmt"

	"github.com/qury53/GC_learn/agent/internal/core/domain"
	core_errors "github.com/qury53/GC_learn/agent/internal/core/errors"
)

const (
	defaultChatsLimit = 50
	maxChatsLimit     = 200
)

func (s *AgentService) GetChats(
	ctx context.Context,
	userID int,
	limit *int,
	offset *int,
) ([]domain.Chat, error) {
	if userID <= 0 {
		return nil, fmt.Errorf(
			"GetChats service: userID has to be positive: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	l := defaultChatsLimit
	if limit != nil {
		l = *limit
	}

	if l <= 0 {
		return nil, fmt.Errorf(
			"GetChats service: limit has to be positive: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if l > maxChatsLimit {
		l = maxChatsLimit
	}

	o := 0
	if offset != nil {
		o = *offset
	}

	if o < 0 {
		return nil, fmt.Errorf(
			"GetChats service: offset cannot be negative: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	chats, err := s.chatRepository.GetChats(ctx, userID, l, o)
	if err != nil {
		return nil, fmt.Errorf("GetChats service: failed to GetChats: %w", err)
	}

	return chats, nil
}
