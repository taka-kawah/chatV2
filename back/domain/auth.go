package domain

import "gorm.io/gorm"

type Auth struct {
	gorm.Model
	Email          string `validate:"required, email" db:"email"`
	HashedPassword string `validate:"required" db:"hashed_password"`
}
