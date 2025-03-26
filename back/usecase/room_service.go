package usecase

import (
	"back/db"
	"back/domain"
	"back/provider"
)

type RoomService struct {
	repo db.RoomDriver
}

func NewRoomService(repo *db.RoomDriver) provider.RoomProvider {
	return &RoomService{repo: *repo}
}

func (rs *RoomService) CreateNewRoom(name string) provider.CustomError {
	if err := rs.repo.Create(name); err != nil {
		return err
	}
	return nil
}

func (rs *RoomService) GetAllRooms() ([]domain.Room, provider.CustomError) {
	rooms, err := rs.repo.FetchAll()
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

func (rs *RoomService) GetRoomById(id uint) (*domain.Room, provider.CustomError) {
	room, err := rs.repo.FetchById(id)
	if err != nil {
		return nil, err
	}
	return room, err
}

func (rs *RoomService) UpdateRoomName(roomId uint, newName string) provider.CustomError {
	if err := rs.repo.UpdateNameById(roomId, newName); err != nil {
		return err
	}
	return nil
}
