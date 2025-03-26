package db

import (
	"errors"
	"log"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/validator/v10"
)

func TestAuthRepo(t *testing.T) {
	mockDbInstances, err := NewMockDbInstances()
	if err != nil {
		log.Fatal(mockDbInstances.GormDb)
	}
	d := NewAuthDriver(mockDbInstances.GormDb)

	t.Run("normal: create auth", func(t *testing.T) {
		mockDbInstances.Mock.ExpectBegin()
		mockDbInstances.Mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "auths" ("created_at","updated_at","deleted_at","email","hashed_password") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "test@test.com", "test_hashed").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mockDbInstances.Mock.ExpectCommit()
		if err := d.Create("test@test.com", "test_hashed"); err != nil {
			t.Errorf("unexpected error (%v)", err)
		}
	})

	t.Run("abnormal: create auth without email", func(t *testing.T) {
		mockDbInstances.Mock.ExpectBegin()
		mockDbInstances.Mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "auths" ("created_at","updated_at","deleted_at","email","hashed_password") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "", "test_hashed")
		mockDbInstances.Mock.ExpectRollback()
		err := d.Create("", "test_hashed")
		if err == nil {
			t.Errorf("expected error but got nil")
			return
		}
		if !errors.As(err.Unwrap(), &validator.ValidationErrors{}) {
			t.Errorf("unexpected error (%v)", err)
			return
		}
	})

	t.Run("abnormal: create auth without password", func(t *testing.T) {
		mockDbInstances.Mock.ExpectBegin()
		mockDbInstances.Mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "auths" ("created_at","updated_at","deleted_at","email","hashed_password") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "test@test.com", "").
			WillReturnError(errors.New("expected"))
		mockDbInstances.Mock.ExpectRollback()
		err := d.Create("test@test.com", "")
		if err == nil {
			t.Errorf("expected error but got nil")
			return
		}
		if err.Unwrap().Error() != "expected" {
			t.Errorf("unexpected error (%v)", err)
			return
		}
	})

	t.Run("normal: fetch auth", func(t *testing.T) {
		mockDbInstances.Mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "auths" WHERE (email = $1 AND hashed_password = $2) AND "auths"."deleted_at" IS NULL ORDER BY "auths"."id" LIMIT $3`)).
			WithArgs("test@test.com", "test_hash", 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "hashed_password"}).AddRow(1, "test@test.com", "test_hash"))

		auth, err := d.CheckIfExist("test@test.com", "test_hash")
		if err != nil {
			t.Errorf("unexpected error (%v)", err)
			return
		}
		if auth == nil {
			t.Errorf("expected auth but got nil")
			return
		}
	})

	t.Run("normal: fetch no auth", func(t *testing.T) {
		mockDbInstances.Mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "auths" WHERE (email = $1 AND hashed_password = $2) AND "auths"."deleted_at" IS NULL ORDER BY "auths"."id" LIMIT $3`)).
			WithArgs("test@test.com", "test_hash", 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "hashed_password"}))

		auth, err := d.CheckIfExist("test@test.com", "test_hash")
		if err != nil {
			t.Errorf("unexpected error (%v)", err)
			return
		}
		if auth != nil {
			t.Errorf("expected nil but got auth %v", auth)
			return
		}
	})

	t.Run("normal: delete auth", func(t *testing.T) {
		mockDbInstances.Mock.ExpectBegin()
		mockDbInstances.Mock.ExpectExec(regexp.QuoteMeta(`UPDATE "auths" SET "deleted_at"=$1 WHERE (email = $2 AND hashed_password = $3) AND "auths"."deleted_at" IS NULL`)).
			WithArgs(sqlmock.AnyArg(), "test@test.com", "test_hashed").
			WillReturnResult(sqlmock.NewResult(int64(1), 1))
		mockDbInstances.Mock.ExpectCommit()

		if err := d.DeleteAuth("test@test.com", "test_hashed"); err != nil {
			t.Errorf("unexpected error (%v)", err)
		}
	})
}
