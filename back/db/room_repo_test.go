package db

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestRoomRepo(t *testing.T) {
	mockDbInstances, err := newMockDbInstances()
	if err != nil {
		log.Fatal("failed to create mock")
	}
	defer mockDbInstances.Disconnect()

	d := NewRoomDriver(mockDbInstances.GormDb)
	t.Run("normal: fetch no room", func(t *testing.T) {
		testFetchAllNone(t, mockDbInstances.Mock, d)
	})
	t.Run("normal: create room", func(t *testing.T) {
		testCreateRoomNormal(t, mockDbInstances.Mock, d, "test1")
	})
	t.Run("normal: fetch 1 room", func(t *testing.T) {
		testFetchAllRooms1Normal(t, mockDbInstances.Mock, d)
	})
	t.Run("normal: fetch 2 rooms", func(t *testing.T) {
		testFetchAllRooms2Normal(t, mockDbInstances.Mock, d)
	})

	t.Run("normal: fetch by id", func(t *testing.T) {
		testFetchRoomById(t, mockDbInstances.Mock, d)
	})

	t.Run("normal: update name by id", func(t *testing.T) {
		testUpdateRoomName(t, mockDbInstances.Mock, d)
	})

	t.Run("normal: delete room by id", func(t *testing.T) {
		testDeleteRoom(t, mockDbInstances.Mock, d)
	})

	t.Run("abNormal: fetch no room by id", func(t *testing.T) {
		testFetchRoomByIdNone(t, mockDbInstances.Mock, d)
	})
	t.Run("abNormal: create room with no name", func(t *testing.T) {
		testCreateRoomWithNoName(t, mockDbInstances.Mock, d)
	})
	t.Run("abnormal: update nonexist room", func(t *testing.T) {
		testUpdateNonExistRoom(t, mockDbInstances.Mock, d)
	})
	t.Run("abnormal: delete nonexist room", func(t *testing.T) {
		testDeleteNonExistRoom(t, mockDbInstances.Mock, d)
	})
}

func testCreateRoomNormal(t *testing.T, m sqlmock.Sqlmock, d *RoomDriver, name string) {
	m.ExpectBegin()
	m.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "rooms" ("created_at","updated_at","deleted_at","name") VALUES ($1,$2,$3,$4) RETURNING "id"`)).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			name,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	m.ExpectCommit()
	if err := d.Create(name); err != nil {
		t.Errorf("unexpected error (%v)", err)
	}
}

func testCreateRoomWithNoName(t *testing.T, m sqlmock.Sqlmock, d *RoomDriver) {
	m.ExpectBegin()
	m.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "rooms" ("created_at","updated_at","deleted_at","name") VALUES ($1,$2,$3,$4) RETURNING "id"`)).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			"",
		).
		WillReturnError(fmt.Errorf("expected"))
	m.ExpectRollback()

	err := d.Create("")
	if err == nil {
		t.Errorf("expected error but got nil")
		return
	}
	if err.Unwrap().Error() != "expected" {
		t.Errorf("unexpected error (%v)", err)
		return
	}
}

func testFetchAllNone(t *testing.T, m sqlmock.Sqlmock, d *RoomDriver) {
	m.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "rooms" WHERE "rooms"."deleted_at" IS NULL`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	rooms, err := d.FetchAll()
	if err != nil {
		t.Errorf("unexpected error (%v)", err)
	}
	if len(rooms) > 0 {
		t.Errorf("expected no rooms but got %v", rooms)
	}
}

func testFetchAllRooms1Normal(t *testing.T, m sqlmock.Sqlmock, d *RoomDriver) {
	m.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "rooms" WHERE "rooms"."deleted_at" IS NULL`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "test"),
		)

	rooms, err := d.FetchAll()
	if err != nil {
		t.Errorf("unexpected error (%v)", err)
		return
	}
	if len(rooms) != 1 {
		t.Errorf("expected 1 rooms but got %v", len(rooms))
		return
	}
}

func testFetchAllRooms2Normal(t *testing.T, m sqlmock.Sqlmock, d *RoomDriver) {
	m.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "rooms" WHERE "rooms"."deleted_at" IS NULL`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "test1").
			AddRow(2, "test2"),
		)

	rooms, err := d.FetchAll()
	if err != nil {
		t.Errorf("unexpected error (%v)", err)
		return
	}
	if len(rooms) != 2 {
		t.Errorf("expected 2 rooms but got %v", len(rooms))
		return
	}
}

func testFetchRoomById(t *testing.T, m sqlmock.Sqlmock, d *RoomDriver) {
	m.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "rooms" WHERE "rooms"."id" = $1 AND "rooms"."deleted_at" IS NULL ORDER BY "rooms"."id" LIMIT $2`)).
		WithArgs(1, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "test"),
	)

	room, err := d.FetchById(1)
	if err != nil {
		t.Errorf("unexpected error (%v)", err)
		return
	}
	if room.ID != 1 {
		t.Errorf("expected id %v but got %v (%v)", 1, room.ID, room)
		return
	}
}

func testFetchRoomByIdNone(t *testing.T, m sqlmock.Sqlmock, d *RoomDriver) {
	id := 1
	m.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "rooms" WHERE "rooms"."id" = $1 AND "rooms"."deleted_at" IS NULL ORDER BY "rooms"."id" LIMIT $2`)).
		WithArgs(id, 1).
		WillReturnError(errors.New("expected"))

	_, err := d.FetchById(uint(id))
	if err == nil {
		t.Errorf("expected error but got nil")
		return
	}
	if err.Unwrap().Error() != "expected" {
		t.Errorf("unexpected error (%v)", err.Unwrap().Error())
		return
	}
}

func testUpdateRoomName(t *testing.T, m sqlmock.Sqlmock, d *RoomDriver) {
	newName := "testupdated"
	m.ExpectBegin()
	m.ExpectExec(regexp.QuoteMeta(`UPDATE "rooms" SET "name"=$1,"updated_at"=$2 WHERE id = $3 AND "rooms"."deleted_at" IS NULL`)).
		WithArgs(newName, sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	m.ExpectCommit()

	if err := d.UpdateNameById(1, newName); err != nil {
		t.Errorf("unexpected error (%v)", err)
		return
	}
}

func testUpdateNonExistRoom(t *testing.T, m sqlmock.Sqlmock, d *RoomDriver) {
	id := 1
	newName := "testupdated"
	m.ExpectBegin()
	m.ExpectExec(regexp.QuoteMeta(`UPDATE "rooms" SET "name"=$1,"updated_at"=$2 WHERE id = $3 AND "rooms"."deleted_at" IS NULL`)).
		WithArgs(newName, sqlmock.AnyArg(), id).
		WillReturnError(fmt.Errorf("expected"))
	m.ExpectRollback()

	err := d.UpdateNameById(uint(id), newName)
	if err == nil {
		t.Errorf("expected error but got nil")
		return
	}
	if err.Unwrap().Error() != "expected" {
		t.Errorf("unexpected error (%v)", err)
		return
	}
}

func testDeleteRoom(t *testing.T, m sqlmock.Sqlmock, d *RoomDriver) {
	id := 1
	m.ExpectBegin()
	m.ExpectExec(regexp.QuoteMeta(`UPDATE "rooms" SET "deleted_at"=$1 WHERE "rooms"."id" = $2 AND "rooms"."deleted_at" IS NULL`)).
		WithArgs(sqlmock.AnyArg(), id).
		WillReturnResult(sqlmock.NewResult(int64(id), 1))
	m.ExpectCommit()

	err := d.DeleteById(uint(id))
	if err != nil {
		t.Errorf("unexpected error (%v)", err)
		return
	}
}

func testDeleteNonExistRoom(t *testing.T, m sqlmock.Sqlmock, d *RoomDriver) {
	id := 1
	m.ExpectBegin()
	m.ExpectExec(regexp.QuoteMeta(`UPDATE "rooms" SET "deleted_at"=$1 WHERE "rooms"."id" = $2 AND "rooms"."deleted_at" IS NULL`)).
		WithArgs(sqlmock.AnyArg(), id).
		WillReturnError(fmt.Errorf("expected"))
	m.ExpectRollback()

	err := d.DeleteById(uint(id))
	if err == nil {
		t.Errorf("expected error but got nil")
		return
	}
	if err.Unwrap().Error() != "expected" {
		t.Errorf("unexpected error (%v)", err)
		return
	}
}
