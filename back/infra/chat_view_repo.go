package infra

import (
	"back/domain"
	"fmt"

	"gorm.io/gorm"
)

type ChatViewDriver struct {
	gormDb gorm.DB
}

func NewChatViewDriver(gormDb *gorm.DB) *ChatViewDriver {
	return &ChatViewDriver{gormDb: *gormDb}
}

func (cvd *ChatViewDriver) FetchRecent(limit int) ([]domain.ChatView, error) {
	var chatViews []domain.ChatView
	query := `SELECT
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
	ORDER BY chats.created_at DESC
	LIMIT ?`
	if err := cvd.gormDb.Raw(query, limit).Scan(&chatViews).Error; err != nil {
		return nil, &ChatViewRepositoryError{msg: "failed to fetch recent chats", err: err}
	}
	return chatViews, nil
}

type ChatViewRepositoryError struct {
	msg string
	err error
}

func (e *ChatViewRepositoryError) Error() string {
	return fmt.Sprintf("error occurs in reading chat_view %s (%s)", e.msg, e.err)
}

func (e *ChatViewRepositoryError) Unwrap() error {
	return e.err
}
