package infra

import (
	"fmt"
	"log"
	"math/rand"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestChatRepo(t *testing.T) {
	mockDbInstances, err := NewMockDbInstances()
	if err != nil {
		log.Fatal("failed to create mock", err)
	}
	defer mockDbInstances.Disconnect()

	d := NewChatDriver(mockDbInstances.GormDb)
	t.Run("normal: create chat", func(t *testing.T) {
		testCreateNormal(t, mockDbInstances.Mock, d, "test")
	})
	t.Run("abnormal: empty value", func(t *testing.T) {
		testAbnormal(t, mockDbInstances.Mock, d, "")
	})
}

func testCreateNormal(t *testing.T, m sqlmock.Sqlmock, d *ChatDriver, value string) {
	m.ExpectBegin()
	m.ExpectQuery(`INSERT INTO "chats"`).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			value,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	m.ExpectCommit()

	if err := d.Create(value, uint(rand.Uint64()), uint(rand.Uint64())); err != nil {
		t.Errorf("unexpected error (%v)", err)
	}
}

func testAbnormal(t *testing.T, m sqlmock.Sqlmock, d *ChatDriver, value string) {
	m.ExpectBegin()
	m.ExpectQuery(`INSERT INTO "chats"`).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			value,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnError(fmt.Errorf("Db error"))
	m.ExpectRollback()
	if err := d.Create(value, uint(rand.Uint64()), uint(rand.Uint64())); err == nil {
		t.Errorf("expected error but got nil (%v)", err)
	}
}
