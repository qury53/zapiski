package agent_service

import (
	"context"
	"fmt"

	"github.com/qury53/GC_learn/agent/internal/domain"
	core_errors "github.com/qury53/GC_learn/agent/internal/errors"
)

const (
	defaultMessagesLimit = 1000
)

func (s *AgentService) GetMessages(
	ctx context.Context,
	chatID int,
	limit *int,
	offset *int,
) ([]domain.Message, error) {
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

	messages, err := s.messageRepository.GetMessages(ctx, chatID, l, o)
	if err != nil {
		return nil, fmt.Errorf("GetMessages service: failed to get messages: %w", err)
	}

	return messages, nil
}
