package usecase

import (
	"back/db"
	"back/domain"
	"back/interfaces"
	"fmt"
)

type RoomService struct {
	repo db.RoomDriver
}

func NewRoomService(repo *db.RoomDriver) interfaces.RoomProvider {
	return &RoomService{repo: *repo}
}

func (rs *RoomService) CreateNewRoom(name string) interfaces.CustomError {
	if err := rs.repo.Create(name); err != nil {
		return &roomServiceError{msg: "failed to create new room", err: err}
	}
	return nil
}

func (rs *RoomService) GetAllRooms() ([]domain.Room, interfaces.CustomError) {
	rooms, err := rs.repo.FetchAll()
	if err != nil {
		return nil, &roomServiceError{msg: "failed to get all rooms", err: err}
	}
	return rooms, nil
}

func (rs *RoomService) UpdateRoomName(roomId uint, newName string) interfaces.CustomError {
	if err := rs.repo.UpdateNameById(roomId, newName); err != nil {
		return &roomServiceError{msg: fmt.Sprintf("failed to update room name: id = %v", roomId), err: err}
	}
	return nil
}

type roomServiceError struct {
	msg string
	err error
}

func (e *roomServiceError) Error() string {
	return fmt.Sprintf("error occurs in room service %s, (%s)", e.msg, e.err)
}

func (e *roomServiceError) Unwrap() error {
	return e.err
}
