package domain

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name           string `validate:"required" db:"name"`
	HashedPasseord string `validate:"required" db:"hashed_password"`
	Email          string `validate:"required, email" db:"email"`
}
