package db

import (
	"errors"
	"log"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestUserRepo(t *testing.T) {
	mockDbInstances, err := NewMockDbInstances()
	if err != nil {
		log.Fatal("failed to create mock")
	}
	d := NewUserDriver(mockDbInstances.GormDb)

	t.Run("normal: create user", func(t *testing.T) {
		testCreateUser(t, mockDbInstances.Mock, d, "test_name", "test@email.com")
	})

	t.Run("abnormal: create user without name", func(t *testing.T) {
		testCreateUserWithoutName(t, mockDbInstances.Mock, d, "test_email")
	})

	t.Run("abnormal: create user without name", func(t *testing.T) {
		testCreateUserWithoutEmail(t, mockDbInstances.Mock, d, "test_name")
	})

	t.Run("normal: fetch user by email", func(t *testing.T) {
		testFetchByEmail(t, mockDbInstances.Mock, d, "test_email")
	})

	t.Run("normal: fetch user by email but none", func(t *testing.T) {
		testFetchByEmailButNone(t, mockDbInstances.Mock, d, "test")
	})

	t.Run("abnormal: fetch user by email without email", func(t *testing.T) {
		testFetchByEmailWithoutEmail(t, mockDbInstances.Mock, d)
	})

	t.Run("normal: fetch all", func(t *testing.T) {
		testFetchAll(t, mockDbInstances.Mock, d)
	})

	t.Run("normal: update name", func(t *testing.T) {
		testUpdateUserById(t, mockDbInstances.Mock, d, 1, "test_updated")
	})

	t.Run("abnormal: update empty name", func(t *testing.T) {
		testUpdateUserWithEmptyName(t, mockDbInstances.Mock, d, 1)
	})

	t.Run("normal: delete by id", func(t *testing.T) {
		testDeleteUser(t, mockDbInstances.Mock, d, 1)
	})
}

func testCreateUser(t *testing.T, m sqlmock.Sqlmock, d *UserDriver, name string, email string) {
	m.ExpectBegin()
	m.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("created_at","updated_at","deleted_at","name","email") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), name, email).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	m.ExpectCommit()
	if err := d.Create(name, email); err != nil {
		t.Errorf("unexpected error (%v)", err)
	}
}

func testCreateUserWithoutName(t *testing.T, m sqlmock.Sqlmock, d *UserDriver, email string) {
	m.ExpectBegin()
	m.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("created_at","updated_at","deleted_at","name","email") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "", email).
		WillReturnError(errors.New("expected"))
	m.ExpectRollback()

	err := d.Create("", email)
	if err == nil {
		t.Errorf("expected error but got nil")
		return
	}
	if err.Unwrap().Error() != "expected" {
		t.Errorf("unecpected error (%v)", err)
		return
	}
}

func testCreateUserWithoutEmail(t *testing.T, m sqlmock.Sqlmock, d *UserDriver, name string) {
	m.ExpectBegin()
	m.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("created_at","updated_at","deleted_at","name","email") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), name, "").
		WillReturnError(errors.New("expected"))
	m.ExpectRollback()

	err := d.Create(name, "")
	if err == nil {
		t.Errorf("expected error but got nil")
		return
	}
	if err.Unwrap().Error() != "expected" {
		t.Errorf("unecpected error (%v)", err)
	}
}

func testFetchByEmail(t *testing.T, m sqlmock.Sqlmock, d *UserDriver, email string) {
	m.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(email, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email"}).AddRow(1, "test", "test_email"))

	user, err := d.FetchByEmail(email)
	if err != nil {
		t.Errorf("unexpected error (%v)", err)
		return
	}
	if user.Email != email {
		t.Errorf("expected email %v but got %v", email, user.Email)
		return
	}
}

func testFetchByEmailButNone(t *testing.T, m sqlmock.Sqlmock, d *UserDriver, email string) {
	m.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(email, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	user, err := d.FetchByEmail(email)
	if err != nil {
		t.Errorf("unexpected error (%v)", err)
		return
	}
	if user != nil {
		t.Errorf("expected nil but got %v", user)
		return
	}
}

func testFetchByEmailWithoutEmail(t *testing.T, m sqlmock.Sqlmock, d *UserDriver) {
	m.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs("", 1).
		WillReturnError(errors.New("expected"))

	_, err := d.FetchByEmail("")
	if err == nil {
		t.Errorf("unexpected error but got nil")
		return
	}
	if err.Unwrap().Error() != "expected" {
		t.Errorf("unexpected error (%v)", err)
	}
}

func testFetchAll(t *testing.T, m sqlmock.Sqlmock, d *UserDriver) {
	m.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."deleted_at" IS NULL`)).
		WithoutArgs().
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	_, err := d.FetchAll()
	if err != nil {
		t.Errorf("unexpected error (%v)", err)
	}
}

func testUpdateUserById(t *testing.T, m sqlmock.Sqlmock, d *UserDriver, id int, newName string) {
	m.ExpectBegin()
	m.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "name"=$1,"updated_at"=$2 WHERE id = $3 AND "users"."deleted_at" IS NULL`)).
		WithArgs(newName, sqlmock.AnyArg(), id).
		WillReturnResult(sqlmock.NewResult(int64(id), 1))
	m.ExpectCommit()

	if err := d.UpdateNameById(uint(id), newName); err != nil {
		t.Errorf("unexpected error (%v)", err)
	}
}

func testUpdateUserWithEmptyName(t *testing.T, m sqlmock.Sqlmock, d *UserDriver, id int) {
	m.ExpectBegin()
	m.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "name"=$1,"updated_at"=$2 WHERE id = $3 AND "users"."deleted_at" IS NULL`)).
		WithArgs("", sqlmock.AnyArg(), id).
		WillReturnError(errors.New("expected"))
	m.ExpectRollback()

	err := d.UpdateNameById(uint(id), "")
	if err == nil {
		t.Errorf("expected error but got nil")
		return
	}
	if err.Unwrap().Error() != "expected" {
		t.Errorf("unexpected error (%v)", err)
		return
	}
}

func testDeleteUser(t *testing.T, m sqlmock.Sqlmock, d *UserDriver, id int) {
	m.ExpectBegin()
	m.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "deleted_at"=$1 WHERE "users"."id" = $2 AND "users"."deleted_at" IS NULL`)).
		WithArgs(sqlmock.AnyArg(), id).
		WillReturnResult(sqlmock.NewResult(int64(id), 1))
	m.ExpectCommit()

	if err := d.DeleteById(uint(id)); err != nil {
		t.Errorf("unexpected error (%v)", err)
	}
}
