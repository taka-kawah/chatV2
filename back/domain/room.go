package domain

import "gorm.io/gorm"

type Room struct {
	gorm.Model
	Name string `validate:"required" db:"name"`
}
