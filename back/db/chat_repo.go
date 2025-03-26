package db

import (
	"back/domain"
	"back/provider"
	"fmt"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type ChatDriver struct {
	gormDb   *gorm.DB
	validate *validator.Validate
}

func NewChatDriver(gormDb *gorm.DB) *ChatDriver {
	return &ChatDriver{gormDb: gormDb, validate: validator.New()}
}

func (cd *ChatDriver) Create(message string, userId uint, roomId uint) provider.CustomError {
	newChat := &domain.Chat{Message: message, UserId: userId, RoomId: roomId}
	if err := cd.validate.Struct(newChat); err != nil {
		return &chatRepositoryError{msg: "validation failure", err: err}
	}
	if err := cd.gormDb.Create(newChat).Error; err != nil {
		return &chatRepositoryError{msg: "failed to create new chat", err: err}
	}
	return nil
}

type chatRepositoryError struct {
	msg string
	err error
}

func (e *chatRepositoryError) Error() string {
	return fmt.Sprintf("error in creating chat db %s (%s)", e.msg, e.err)
}

func (e *chatRepositoryError) Unwrap() error {
	return e.err
}
