package db

import (
	"back/domain"
	"back/provider"
	"fmt"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type UserDriver struct {
	gormDb   *gorm.DB
	validate *validator.Validate
}

func NewUserDriver(gormDb *gorm.DB) *UserDriver {
	return &UserDriver{gormDb: gormDb}
}

func (ud *UserDriver) Create(name string, email string) provider.CustomError {
	newUser := domain.User{Name: name, Email: email}
	if err := ud.validate.Struct(&newUser); err != nil {
		return &userRepositoryError{msg: "validation failure", err: err}
	}
	if err := ud.gormDb.Create(&newUser).Error; err != nil {
		return &userRepositoryError{msg: "failed to create new user", err: err}
	}
	return nil
}

func (ud *UserDriver) FetchByEmail(email string) (*domain.User, provider.CustomError) {
	var user domain.User
	res := ud.gormDb.Where("email = ?", email).First(&user)
	if res.RowsAffected == 0 {
		return nil, nil
	}
	if res.Error != nil {
		return nil, &userRepositoryError{msg: fmt.Sprintf("failed to fetch user: email = %v", email), err: res.Error}
	}
	return &user, nil
}

func (ud *UserDriver) FetchAll() (*[]domain.User, provider.CustomError) {
	var users []domain.User
	res := ud.gormDb.Find(&users)
	if res.Error != nil {
		return nil, &userRepositoryError{msg: "failed to fetch all users", err: res.Error}
	}
	return &users, nil
}

func (ud *UserDriver) UpdateNameById(id uint, newName string) provider.CustomError {
	if err := ud.gormDb.Model(&domain.User{}).Where("id = ?", id).Update("name", newName).Error; err != nil {
		return &userRepositoryError{msg: fmt.Sprintf("failed to update user: id = %v", id), err: err}
	}
	return nil
}

func (ud *UserDriver) DeleteById(id uint) provider.CustomError {
	if err := ud.gormDb.Delete(&domain.User{}, id).Error; err != nil {
		return &userRepositoryError{msg: fmt.Sprintf("failed to update user: id = %v", id), err: err}
	}
	return nil
}

type userRepositoryError struct {
	msg string
	err error
}

func (e *userRepositoryError) Error() string {
	return fmt.Sprintf("error in crud user db %s (%s)", e.msg, e.err)
}

func (e *userRepositoryError) Unwrap() error {
	return e.err
}
