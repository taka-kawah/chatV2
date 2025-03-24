package db

import (
	"back/domain"
	"back/interfaces"
	"fmt"

	"gorm.io/gorm"
)

type ChatViewDriver struct {
	gormDb gorm.DB
}

func NewChatViewDriver(gormDb *gorm.DB) *ChatViewDriver {
	return &ChatViewDriver{gormDb: *gormDb}
}

func (cvd *ChatViewDriver) FetchRecent(roomId uint, limit int) ([]domain.ChatView, interfaces.CustomError) {
	var chatViews []domain.ChatView
	query := fmt.Sprintf(`SELECT
	chats.id, 
	chats.created_at, 
	chats.updated_at, 
	chats.deleted_at, 
	chats.message, 
	chats.user_id, 
	chats.room_id, 
	users.name
	FROM chats
	JOIN users ON chats.user_id = users.id
	WHERE chats.room_id = %v
	ORDER BY chats.created_at DESC
	LIMIT %v`, roomId, limit)
	if err := cvd.gormDb.Raw(query).Scan(&chatViews).Error; err != nil {
		return nil, &chatViewRepositoryError{msg: "failed to fetch recent chats", err: err}
	}
	return chatViews, nil
}

type chatViewRepositoryError struct {
	msg string
	err error
}

func (e *chatViewRepositoryError) Error() string {
	return fmt.Sprintf("error occurs in reading chat_view %s (%s)", e.msg, e.err)
}

func (e *chatViewRepositoryError) Unwrap() error {
	return e.err
}
