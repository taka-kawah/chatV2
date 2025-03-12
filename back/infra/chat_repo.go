package infra

import (
	"back/domain"
	"fmt"

	"gorm.io/gorm"
)

type ChatDriver struct {
	gormDb *gorm.DB
}

func NewChatRepository(gormDb *gorm.DB) *ChatDriver {
	return &ChatDriver{gormDb: gormDb}
}

func (cd *ChatDriver) Create(message string, userId uint, roomId uint) *ChatRepositoryError {
	newChat := &domain.Chat{Message: message, UserId: userId, RoomId: roomId}
	if err := cd.gormDb.Create(newChat).Error; err != nil {
		return &ChatRepositoryError{msg: "failed to create new chat", err: err}
	}
	return nil
}

type ChatRepositoryError struct {
	msg string
	err error
}

func (e *ChatRepositoryError) Error() string {
	return fmt.Sprintf("error in creating chat db %s (%s)", e.msg, e.err)
}

func (e *ChatRepositoryError) Unwrap() error {
	return e.err
}
