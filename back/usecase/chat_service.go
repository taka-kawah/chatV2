package usecase

import (
	"back/db"
	"back/domain"
	"fmt"
)

type ChatService struct {
	chatRepo     db.ChatDriver
	chatViewRepo db.ChatViewDriver
}

type IChatService interface {
	PostChat(message string, userId uint, roomId uint) *ChatServiceError
	GetRecentChatsFromOneRoom(roomId uint) (*[]domain.ChatView, *ChatServiceError)
}

func NewChatService(chatRepo *db.ChatDriver, chatViewRepo *db.ChatViewDriver) *ChatService {
	return &ChatService{chatRepo: *chatRepo, chatViewRepo: *chatViewRepo}
}

func (cs *ChatService) PostChat(message string, userId uint, roomId uint) *ChatServiceError {
	if err := cs.chatRepo.Create(message, userId, roomId); err != nil {
		return &ChatServiceError{msg: "failed to post chat", err: err}
	}
	return nil
}

func (cs *ChatService) GetRecentChatsFromOneRoom(roomId uint) (*[]domain.ChatView, *ChatServiceError) {
	limit := 10 //とりあえず直近10件
	chats, err := cs.chatViewRepo.FetchRecent(roomId, limit)
	if err != nil {
		return nil, &ChatServiceError{msg: "failed to read chats", err: err}
	}
	return &chats, nil
}

type ChatServiceError struct {
	msg string
	err error
}

func (e *ChatServiceError) Error() string {
	return fmt.Sprintf("error occurs in chat service %s (%s)", e.msg, e.err)
}

func (e *ChatServiceError) Unwrap() error {
	return e.err
}
