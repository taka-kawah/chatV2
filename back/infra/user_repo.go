package infra

import (
	"back/domain"
	"fmt"

	"gorm.io/gorm"
)

type UserDriver struct {
	gormDb *gorm.DB
}

func NewUserDriver(gormDb *gorm.DB) *UserDriver {
	return &UserDriver{gormDb: gormDb}
}

func (ud *UserDriver) Create(name string, email string, hashedPassword string) error {
	newUser := domain.User{Name: name, Email: email, HashedPasseord: hashedPassword}
	if err := ud.gormDb.Create(&newUser).Error; err != nil {
		return &UserRepositoryError{msg: "failed to create new user", err: err}
	}
	return nil
}

func (ud *UserDriver) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	res := ud.gormDb.First(&user)
	if res.Error != nil {
		return nil, &UserRepositoryError{msg: fmt.Sprintf("failed to find user: email = %v", email), err: res.Error}
	}
	return &user, nil
}

func (ud *UserDriver) FindAll() (*[]domain.User, error) {
	var users []domain.User
	res := ud.gormDb.Find(&users)
	if res.Error != nil {
		return nil, &UserRepositoryError{msg: "failed to find all users", err: res.Error}
	}
	return &users, nil
}

func (ud *UserDriver) Update(user *domain.User) error {
	if err := ud.gormDb.Save(user).Error; err != nil {
		return &UserRepositoryError{msg: fmt.Sprintf("failed to update user: id = %v", user.ID), err: err}
	}
	return nil
}

func (ud *UserDriver) Delete(user *domain.User) error {
	if err := ud.gormDb.Delete(user).Error; err != nil {
		return &UserRepositoryError{msg: fmt.Sprintf("failed to update user: id = %v", user.ID), err: err}
	}
	return nil
}

type UserRepositoryError struct {
	msg string
	err error
}

func (e *UserRepositoryError) Error() string {
	return fmt.Sprintf("error in crud user db %s (%s)", e.msg, e.err)
}

func (e *UserRepositoryError) Unwrap() error {
	return e.err
}
