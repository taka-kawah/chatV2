package usecase

import (
	"back/domain"
	"back/infra"
	"fmt"
)

type RoomService struct {
	repo infra.RoomDriver
}

type IRoomService interface {
	CreateNewRoom(name string) *RoomServiceError
	GetAllRooms() ([]domain.Room, *RoomServiceError)
	UpdateRoomName(roomId uint, newName string) *RoomServiceError
}

func NewRoomService(repo *infra.RoomDriver) *RoomService {
	return &RoomService{repo: *repo}
}

func (rs *RoomService) CreateNewRoom(name string) *RoomServiceError {
	if err := rs.repo.Create(name); err != nil {
		return &RoomServiceError{msg: "failed to create new room", err: err}
	}
	return nil
}

func (rs *RoomService) GetAllRooms() ([]domain.Room, *RoomServiceError) {
	rooms, err := rs.repo.FetchAll()
	if err != nil {
		return nil, &RoomServiceError{msg: "failed to get all rooms", err: err}
	}
	return rooms, nil
}

func (rs *RoomService) UpdateRoomName(roomId uint, newName string) *RoomServiceError {
	if err := rs.repo.UpdateNameById(roomId, newName); err != nil {
		return &RoomServiceError{msg: fmt.Sprintf("failed to update room name: id = %v", roomId), err: err}
	}
	return nil
}

type RoomServiceError struct {
	msg string
	err error
}

func (e *RoomServiceError) Error() string {
	return fmt.Sprintf("error occurs in room service %s, (%s)", e.msg, e.err)
}

func (e *RoomServiceError) Unwrap() error {
	return e.err
}
