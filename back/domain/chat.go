package domain

import (
	"gorm.io/gorm"
)

type Chat struct {
	gorm.Model
	Message string `validate:"required" db:"message"`
	UserId  uint   `validate:"required" db:"user_id"`
	RoomId  uint   `validate:"required" db:"room_id"`
}
