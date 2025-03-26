package usecase

import (
	"back/db"
	"back/domain"
	"back/provider"
	"fmt"
)

type UserService struct {
	repo db.UserDriver
}

func NewUserService(repo *db.UserDriver) provider.UserProvider {
	return &UserService{repo: *repo}
}

func (us *UserService) GetAllUsers() ([]domain.User, provider.CustomError) {
	users, err := us.repo.FetchAll()
	if err != nil {
		return nil, &userServiceError{msg: "failed to get all users", err: err}
	}
	return *users, nil
}

func (us *UserService) GetFromEmail(email string) (*domain.User, provider.CustomError) {
	user, err := us.repo.FetchByEmail(email)
	if err != nil {
		return nil, &userServiceError{msg: fmt.Sprintf("failed to get user from email: %v", email), err: err}
	}
	return user, nil
}

func (us *UserService) RegisterAccount(name string, email string) provider.CustomError {
	if err := us.repo.Create(name, email); err != nil {
		return &userServiceError{msg: "failed to create account", err: err}
	}
	return nil
}

func (us *UserService) UpdateName(id uint, newName string) provider.CustomError {
	if err := us.repo.UpdateNameById(id, newName); err != nil {
		return &userServiceError{msg: fmt.Sprintf("failed to update account id: %v", id), err: err}
	}
	return nil
}

func (us *UserService) Delete(id uint) provider.CustomError {
	if err := us.repo.DeleteById(id); err != nil {
		return &userServiceError{msg: fmt.Sprintf("failed to delete account id: %v", id), err: err}
	}
	return nil
}

type userServiceError struct {
	msg string
	err error
}

func (e *userServiceError) Error() string {
	return fmt.Sprintf("error occurs in user service %s (%s)", e.msg, e.err)
}

func (e *userServiceError) Unwrap() error {
	return e.err
}
