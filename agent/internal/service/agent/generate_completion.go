package agent_service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/qury53/GC_learn/agent/internal/core/domain"
	core_errors "github.com/qury53/GC_learn/agent/internal/core/errors"
)

const (
	// лимит контекстного окна, чем больше, тем пизже ответы в долгосроке, но дороже по токенам
	chatContextLimit = 20
)

func (s *AgentService) GenerateCompletion(
	ctx context.Context,
	request domain.ChatRequest,
) (domain.ChatResponse, error) {
	if request.UserID <= 0 {
		return domain.ChatResponse{}, fmt.Errorf(
			"GenerateCompletion service: userID has to be positive: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	// Роли system/assistant от клиента подменяют системный промпт и историю.
	for _, m := range request.Messages {
		if m.Role != domain.RoleUser {
			return domain.ChatResponse{}, fmt.Errorf(
				"GenerateCompletion service: only user messages are allowed: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	var (
		chat domain.Chat
		err  error
	)

	if request.ChatID == nil {
		chat, err = s.chatRepository.CreateChat(ctx, domain.NewChatUninitialized(
			request.UserID,
			request.ChatMode,
		))

		if err != nil {
			s.log.Info(
				"failed to create chat",
				slog.Int("userID", request.UserID),
				slog.Any("error", err),
			)

			return domain.ChatResponse{}, fmt.Errorf(
				"GenerateCompletion service: failed to create chat: %w",
				err,
			)
		}

		request.ChatID = &chat.ID
	} else {
		chat, err = s.chatRepository.GetChat(ctx, *request.ChatID)
		if err != nil {
			s.log.Info(
				"failed to get chat",
				slog.Int("chatID", *request.ChatID),
				slog.Any("error", err),
			)

			return domain.ChatResponse{}, fmt.Errorf(
				"GenerateCompletion service: failed to get chat: %w",
				err,
			)
		}

		// Чужой чат неотличим от отсутствующего.
		if chat.UserID != request.UserID {
			return domain.ChatResponse{}, fmt.Errorf(
				"GenerateCompletion service: chat not found: %w",
				core_errors.ErrNotFound,
			)
		}

		// Режим берём из БД, а не из запроса клиента.
		request.ChatMode = chat.Mode
	}

	msgContext, err := s.messageRepository.GetMessages(
		ctx,
		*request.ChatID,
		chatContextLimit,
		0,
	)
	if err != nil {
		s.log.Debug(
			"failed to get chat context",
			slog.Int("chatID", *request.ChatID),
			slog.Any("error", err),
		)

		return domain.ChatResponse{}, fmt.Errorf("generate completion: %w", err)
	}

	// Новые сообщения запоминаем до склейки с историей: в склейке к ним
	// добавляется system-промпт, и индекс «сколько было контекста» уже не работает.
	newMessages := request.Messages

	messages := make([]domain.Message, 0, len(msgContext)+len(newMessages)+1)

	if request.ChatMode == domain.ModeSupport {
		messages = append(messages, domain.Message{
			Role:    domain.RoleSystem,
			Content: supportSystemPrompt,
		})
	}

	messages = append(messages, msgContext...)
	messages = append(messages, newMessages...)

	request.Messages = messages
	resp, err := s.aiProvider.GenerateCompletion(ctx, request)
	if err != nil {
		return domain.ChatResponse{}, err
	}

	for _, m := range newMessages {
		if m.Role == domain.RoleUser {
			if _, err := s.messageRepository.Create(ctx, domain.Message{
				ChatID:  *request.ChatID,
				Role:    m.Role,
				Content: m.Content,
			}); err != nil {
				s.log.Error("failed to save user message", slog.Any("error", err))
			}
		}
	}

	if _, err := s.messageRepository.Create(ctx, domain.Message{
		ChatID:  *request.ChatID,
		Role:    domain.RoleAssistant,
		Content: resp.Content,
	}); err != nil {
		s.log.Error("failed to save assistant message", slog.Any("error", err))
	}

	return resp, nil
}
