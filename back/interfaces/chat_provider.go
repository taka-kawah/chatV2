package interfaces

import "back/domain"

type ChatProvider interface {
	PostChat(message string, userId uint, roomId uint) CustomError
	GetRecentChatsFromOneRoom(roomId uint) (*[]domain.ChatView, CustomError)
}
