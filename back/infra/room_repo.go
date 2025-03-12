package infra

import (
	"back/domain"
	"fmt"

	"gorm.io/gorm"
)

type RoomDriver struct {
	gormDb *gorm.DB
}

func NewRoomDriver(gormDb *gorm.DB) *RoomDriver {
	return &RoomDriver{gormDb: gormDb}
}

func (rd *RoomDriver) Create(name string) *RoomRepositoryError {
	newRoom := domain.Room{Name: name}
	if err := rd.gormDb.Create(newRoom).Error; err != nil {
		return &RoomRepositoryError{msg: "failed to create room", err: err}
	}
	return nil
}

func (rd *RoomDriver) FetchAll() (*[]domain.Room, *RoomRepositoryError) {
	var rooms []domain.Room
	res := rd.gormDb.Find(&rooms)
	if res.Error != nil {
		return nil, &RoomRepositoryError{msg: "failed to fetch all rooms", err: res.Error}
	}
	return &rooms, nil
}

func (rd *RoomDriver) Update(room *domain.Room) *RoomRepositoryError {
	if err := rd.gormDb.Save(room).Error; err != nil {
		return &RoomRepositoryError{msg: fmt.Sprintf("failed to update room: id = %v", room.ID), err: err}
	}
	return nil
}

func (rd *RoomDriver) Delete(room *domain.Room) *RoomRepositoryError {
	if err := rd.gormDb.Delete(room).Error; err != nil {
		return &RoomRepositoryError{msg: fmt.Sprintf("failed to delete room: id = %v", room.ID), err: err}
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
