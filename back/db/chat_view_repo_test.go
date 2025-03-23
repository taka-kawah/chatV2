package db

import (
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestChatViewRepo(t *testing.T) {
	mockDbInstances, err := NewMockDbInstances()
	if err != nil {
		log.Fatal("failed to create mock")
	}
	d := NewChatViewDriver(mockDbInstances.GormDb)

	t.Run("normal: fetch view", func(t *testing.T) {
		mockDbInstances.Mock.ExpectQuery(`SELECT
	chats.id, 
	chats.created_at, 
	chats.updated_at, 
	chats.deleted_at, 
	chats.message, 
	chats.user_id, 
	chats.room_id, 
	users.name
	FROM chats
	JOIN users ON chats.user_id = users.id
	WHERE chats.room_id = 1
	ORDER BY chats.created_at DESC
	LIMIT 10`).
			WithoutArgs().
			WillReturnRows(sqlmock.NewRows([]string{`"chats"."id"`, `"users"."name"`}))
	})
	_, err = d.FetchRecent(1, 10)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}
