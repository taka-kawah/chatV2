package usecase

import (
	"back/db"
	"back/domain"
	"fmt"
)

type UserService struct {
	repo db.UserDriver
}

type IUserService interface {
	GetAllUsers() ([]domain.User, *UserServiceError)
	GetFromEmail(email string) (*domain.User, *UserServiceError)
	RegisterAccount(name string, email string) *UserServiceError
	UpdateName(id uint, newName string) *UserServiceError
	Delete(id uint) *UserServiceError
}

func NewUserService(repo *db.UserDriver) *UserService {
	return &UserService{repo: *repo}
}

func (us *UserService) GetAllUsers() ([]domain.User, *UserServiceError) {
	users, err := us.repo.FetchAll()
	if err != nil {
		return nil, &UserServiceError{msg: "failed to get all users", err: err}
	}
	return *users, nil
}

func (us *UserService) GetFromEmail(email string) (*domain.User, *UserServiceError) {
	user, err := us.repo.FetchByEmail(email)
	if err != nil {
		return nil, &UserServiceError{msg: fmt.Sprintf("failed to get user from email: %v", email), err: err}
	}
	return user, nil
}

func (us *UserService) RegisterAccount(name string, email string) *UserServiceError {
	if err := us.repo.Create(name, email); err != nil {
		return &UserServiceError{msg: "failed to create account", err: err}
	}
	return nil
}

func (us *UserService) UpdateName(id uint, newName string) *UserServiceError {
	if err := us.repo.UpdateNameById(id, newName); err != nil {
		return &UserServiceError{msg: fmt.Sprintf("failed to update account id: %v", id), err: err}
	}
	return nil
}

func (us *UserService) Delete(id uint) *UserServiceError {
	if err := us.repo.DeleteById(id); err != nil {
		return &UserServiceError{msg: fmt.Sprintf("failed to delete account id: %v", id), err: err}
	}
	return nil
}

type UserServiceError struct {
	msg string
	err error
}

func (e *UserServiceError) Error() string {
	return fmt.Sprintf("error occurs in user service %s (%s)", e.msg, e.err)
}

func (e *UserServiceError) Unwrap() error {
	return e.err
}
