package db

import (
	"back/domain"
	"back/provider"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type AuthDriver struct {
	gormDb   *gorm.DB
	validate *validator.Validate
}

func NewAuthDriver(gormDb *gorm.DB) *AuthDriver {
	return &AuthDriver{gormDb: gormDb, validate: validator.New()}
}

func (d *AuthDriver) Create(email string, hashedPassword string) provider.CustomError {
	newAuth := &domain.Auth{Email: email, HashedPassword: hashedPassword}
	if err := d.validate.Struct(newAuth); err != nil {
		return &authRepositoryError{msg: "validation failure", err: err}
	}
	if err := d.gormDb.Create(newAuth).Error; err != nil {
		return &authRepositoryError{msg: "failed to create auth", err: err}
	}
	return nil
}

func (d *AuthDriver) CheckIfExist(email string, hashedPassword string) (*domain.Auth, provider.CustomError) {
	var auth domain.Auth
	res := d.gormDb.Where("email = ? AND hashed_password = ?", email, hashedPassword).First(&auth)
	if res.Error == nil {
		return &auth, nil
	}
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return nil, &authRepositoryError{msg: fmt.Sprintf("failled to fetch auth email: %v", email), err: res.Error}
}

func (d *AuthDriver) DeleteAuth(email string, hashedPassword string) provider.CustomError {
	var auth domain.Auth
	if err := d.gormDb.Where("email = ? AND hashed_password = ?", email, hashedPassword).Delete(&auth).Error; err != nil {
		return &authRepositoryError{msg: "failed to delete auth", err: err}
	}
	return nil
}

type authRepositoryError struct {
	msg string
	err error
}

func (e *authRepositoryError) Error() string {
	return fmt.Sprintf("error occurs in auth db %s (%s)", e.msg, e.err)
}

func (e *authRepositoryError) Unwrap() error {
	return e.err
}
