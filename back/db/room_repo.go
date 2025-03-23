package db

import (
	"back/domain"
	"back/interfaces"
	"fmt"

	"gorm.io/gorm"
)

type RoomDriver struct {
	gormDb *gorm.DB
}

func NewRoomDriver(gormDb *gorm.DB) *RoomDriver {
	return &RoomDriver{gormDb: gormDb}
}

func (rd *RoomDriver) Create(name string) interfaces.CustomError {
	newRoom := domain.Room{Name: name}
	if err := rd.gormDb.Create(&newRoom).Error; err != nil {
		return &RoomRepositoryError{msg: "failed to create room", err: err}
	}
	return nil
}

func (rd *RoomDriver) FetchAll() ([]domain.Room, interfaces.CustomError) {
	var rooms []domain.Room
	res := rd.gormDb.Find(&rooms)
	if res.Error != nil {
		return nil, &RoomRepositoryError{msg: "failed to fetch all rooms", err: res.Error}
	}
	return rooms, nil
}

func (rd *RoomDriver) FetchById(id uint) (*domain.Room, interfaces.CustomError) {
	var room domain.Room
	res := rd.gormDb.First(&room, id)
	if res.Error != nil {
		return nil, &RoomRepositoryError{msg: fmt.Sprintf("failed to get room: id = %v", id), err: res.Error}
	}
	return &room, nil
}

func (rd *RoomDriver) UpdateNameById(id uint, newName string) interfaces.CustomError {
	if err := rd.gormDb.Model(&domain.Room{}).Where("id = ?", id).Update("name", newName).Error; err != nil {
		return &RoomRepositoryError{msg: fmt.Sprintf("failed to update room: id = %v", id), err: err}
	}
	return nil
}

func (rd *RoomDriver) DeleteById(id uint) interfaces.CustomError {
	if err := rd.gormDb.Delete(&domain.Room{}, id).Error; err != nil {
		return &RoomRepositoryError{msg: fmt.Sprintf("failed to delete room: id = %v", id), err: err}
	}
	return nil
}

type RoomRepositoryError struct {
	msg string
	err error
}

func (e *RoomRepositoryError) Error() string {
	return fmt.Sprintf("error in crud room db %s (%s)", e.msg, e.err)
}

func (e *RoomRepositoryError) Unwrap() error {
	return e.err
}
