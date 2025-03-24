package interfaces

import (
	"back/domain"
)

type RoomProvider interface {
	CreateNewRoom(name string) CustomError
	GetAllRooms() ([]domain.Room, CustomError)
	UpdateRoomName(roomId uint, newName string) CustomError
}
