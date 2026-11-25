package agent_service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/qury53/GC_learn/agent/internal/domain"
)

const (
	//лимит контекстного окна, чем больше, тем пизже ответы в долгосроке, но дороже по токенам
	chatContextLimit = 20
)

func (s *AgentService) GenerateCompletion(
	ctx context.Context,
	request domain.ChatRequest,
) (domain.ChatResponse, error) {
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

	messages := make([]domain.Message, 0, len(msgContext)+2)
	switch request.ChatMode {
	case domain.ModeSupport:
		messages = append(messages, domain.Message{
			Role:    domain.RoleSystem,
			Content: supportSystemPrompt,
		})
	}

	messages = append(messages, msgContext...)
	messages = append(messages, request.Messages...)

	request.Messages = messages
	resp, err := s.aiProvider.GenerateCompletion(ctx, request)
	if err != nil {
		return domain.ChatResponse{}, err
	}

	for _, m := range request.Messages[len(msgContext)+1:] { // только новые, без system и history
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
