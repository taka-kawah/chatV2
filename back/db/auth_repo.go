package db

import (
	"back/domain"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type AuthDriver struct {
	gormDb *gorm.DB
}

func NewAuthDriver(gormDb *gorm.DB) *AuthDriver {
	return &AuthDriver{gormDb: gormDb}
}

func (d *AuthDriver) Create(email string, hashedPassword string) error {
	newAuth := &domain.Auth{Email: email, HashedPassword: hashedPassword}
	if err := d.gormDb.Create(newAuth).Error; err != nil {
		return &AuthRepositoryError{msg: "failed to create auth", err: err}
	}
	return nil
}

func (d *AuthDriver) CheckIfExist(email string, hashedPassword string) (*domain.Auth, error) {
	var auth domain.Auth
	res := d.gormDb.Where("email = ? AND hashed_password = ?", email, hashedPassword).First(&auth)
	if res.Error == nil {
		return &auth, nil
	}
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return nil, &AuthRepositoryError{msg: fmt.Sprintf("failled to fetch auth email: %v", email), err: res.Error}
}

func (d *AuthDriver) DeleteAuth(email string, hashedPassword string) error {
	var auth domain.Auth
	if err := d.gormDb.Where("email = ? AND hashed_password = ?", email, hashedPassword).Delete(&auth).Error; err != nil {
		return &AuthRepositoryError{msg: "failed to delete auth", err: err}
	}
	return nil
}

type AuthRepositoryError struct {
	msg string
	err error
}

func (e *AuthRepositoryError) Error() string {
	return fmt.Sprintf("error occurs in auth db %s (%s)", e.msg, e.err)
}

func (e *AuthRepositoryError) Unwrap() error {
	return e.err
}
