package domain

import "gorm.io/gorm"

type Room struct {
	gorm.Model
	Name string `validate:"required" db:"name"`
}

func (r *Room) GetNameUpdated(newName string) *Room {
	newRoom := r
	newRoom.Name = newName
	return newRoom
}
