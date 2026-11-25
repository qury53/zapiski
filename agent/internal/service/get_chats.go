package agent_service

import (
	"context"
	"fmt"

	"github.com/qury53/GC_learn/agent/internal/domain"
	core_errors "github.com/qury53/GC_learn/agent/internal/errors"
)

func (s *AgentService) GetChats(
	ctx context.Context,
	userID int,
	limit *int,
	offset *int,
) ([]domain.Chat, error) {
	if limit != nil && *limit <= 0 {
		return nil, fmt.Errorf(
			"GetChats service: limit has to be positive: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf(
			"GetChats service: offset cannot be negative: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	chats, err := s.chatRepository.GetChats(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("GetChats service: failed to GetChats: %w", err)
	}

	return chats, nil
}
