package domain

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name  string `validate:"required" db:"name"`
	Email string `validate:"required, email" db:"email"`
}
