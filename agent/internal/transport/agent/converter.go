package agent_transport

import (
	v1 "shared/pkg/proto/agent/v1"

	"github.com/qury53/GC_learn/agent/internal/core/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToGetChatsResponse конвертирует слайс доменных сущностей в protobuf-ответ.
func ToGetChatsResponse(chats []domain.Chat) *v1.GetChatsResponse {
	pbChats := make([]*v1.Chat, 0, len(chats))
	for _, chat := range chats {
		pbChats = append(pbChats, ToProtoChat(chat))
	}

	return &v1.GetChatsResponse{
		Chats: pbChats,
	}
}

// ToProtoChat конвертирует одну доменную сущность Chat в protobuf-сообщение.
func ToProtoChat(c domain.Chat) *v1.Chat {
	pbChat := &v1.Chat{
		Id:        int64(c.ID),
		Version:   int64(c.Version),
		Mode:      toProtoChatMode(c.Mode),
		CreatedAt: timestamppb.New(c.CreatedAt),
		UpdatedAt: timestamppb.New(c.UpdatedAt),
		UserId:    int64(c.UserID),
	}

	if c.DeletedAt != nil {
		pbChat.DeletedAt = timestamppb.New(*c.DeletedAt)
	}

	return pbChat
}

// toProtoChatMode маппит доменный enum в protobuf enum.
func toProtoChatMode(mode domain.ChatMode) v1.ChatMode {
	switch mode {
	case domain.ModeSupport:
		return v1.ChatMode_CHAT_MODE_SUPPORT
	case domain.ModeFull:
		return v1.ChatMode_CHAT_MODE_FULL
	default:
		return v1.ChatMode_CHAT_MODE_UNSPECIFIED
	}
}

// ToGetMessagesResponse конвертирует слайс доменных сообщений в protobuf-ответ.
func ToGetMessagesResponse(messages []domain.Message) *v1.GetMessagesResponse {
	pbMessages := make([]*v1.Message, 0, len(messages))
	for _, msg := range messages {
		pbMessages = append(pbMessages, ToProtoMessage(msg))
	}

	return &v1.GetMessagesResponse{
		Message: pbMessages,
	}
}

// ToProtoMessage конвертирует одно доменное сообщение в protobuf-сообщение.
func ToProtoMessage(m domain.Message) *v1.Message {
	return &v1.Message{
		Id:        int64(m.ID),
		ChatId:    int64(m.ChatID),
		Role:      toProtoRole(m.Role),
		Content:   m.Content,
		CreatedAt: timestamppb.New(m.CreatedAt),
	}
}

// toProtoRole маппит доменную роль в protobuf enum pb.Role.
// Проверьте названия констант из вашего пакета domain (RoleSystem, RoleUser и т.д.).
func toProtoRole(role domain.Role) v1.Role {
	switch role {
	case domain.RoleSystem:
		return v1.Role_ROLE_SYSTEM
	case domain.RoleUser:
		return v1.Role_ROLE_USER
	case domain.RoleAssistant:
		return v1.Role_ROLE_ASSISTANT
	case domain.RoleTool:
		return v1.Role_ROLE_TOOL
	default:
		return v1.Role_ROLE_UNSPECIFIED
	}
}

func int64PtrToIntPtr(p *int64) *int {
	if p == nil {
		return nil
	}
	v := int(*p)
	return &v
}
