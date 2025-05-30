package provider

import (
	"back/domain"
)

type UserProvider interface {
	GetAllUsers() ([]domain.User, CustomError)
	GetFromEmail(email string) (*domain.User, CustomError)
	RegisterAccount(name string, email string) CustomError
	UpdateName(id uint, newName string) CustomError
	Delete(id uint) CustomError
}
