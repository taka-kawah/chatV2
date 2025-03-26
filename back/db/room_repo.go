package db

import (
	"back/domain"
	"back/provider"
	"fmt"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type RoomDriver struct {
	gormDb   *gorm.DB
	validate *validator.Validate
}

func NewRoomDriver(gormDb *gorm.DB) *RoomDriver {
	return &RoomDriver{gormDb: gormDb, validate: validator.New()}
}

func (rd *RoomDriver) Create(name string) provider.CustomError {
	newRoom := domain.Room{Name: name}
	if err := rd.validate.Struct(&newRoom); err != nil {
		return &roomRepositoryError{msg: "validation failure", err: err}
	}
	if err := rd.gormDb.Create(&newRoom).Error; err != nil {
		return &roomRepositoryError{msg: "failed to create room", err: err}
	}
	return nil
}

func (rd *RoomDriver) FetchAll() ([]domain.Room, provider.CustomError) {
	var rooms []domain.Room
	res := rd.gormDb.Find(&rooms)
	if res.Error != nil {
		return nil, &roomRepositoryError{msg: "failed to fetch all rooms", err: res.Error}
	}
	return rooms, nil
}

func (rd *RoomDriver) FetchById(id uint) (*domain.Room, provider.CustomError) {
	var room domain.Room
	res := rd.gormDb.First(&room, id)
	if res.Error != nil {
		return nil, &roomRepositoryError{msg: fmt.Sprintf("failed to get room: id = %v", id), err: res.Error}
	}
	return &room, nil
}

func (rd *RoomDriver) UpdateNameById(id uint, newName string) provider.CustomError {
	if err := rd.gormDb.Model(&domain.Room{}).Where("id = ?", id).Update("name", newName).Error; err != nil {
		return &roomRepositoryError{msg: fmt.Sprintf("failed to update room: id = %v", id), err: err}
	}
	return nil
}

func (rd *RoomDriver) DeleteById(id uint) provider.CustomError {
	if err := rd.gormDb.Delete(&domain.Room{}, id).Error; err != nil {
		return &roomRepositoryError{msg: fmt.Sprintf("failed to delete room: id = %v", id), err: err}
	}
	return nil
}

type roomRepositoryError struct {
	msg string
	err error
}

func (e *roomRepositoryError) Error() string {
	return fmt.Sprintf("error in crud room db %s (%s)", e.msg, e.err)
}

func (e *roomRepositoryError) Unwrap() error {
	return e.err
}
