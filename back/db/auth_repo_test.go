package db

import (
	"errors"
	"log"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/validator/v10"
)

func setUpAuthDriver() (*mockDbInstances, *AuthDriver) {
	mock, err := newMockDbInstances()
	if err != nil {
		log.Fatal("failed to create mock", err)
	}
	d := NewAuthDriver(mock.GormDb)
	return mock, d
}
func TestAuthRepo(t *testing.T) {

	t.Run("normal: create auth", func(t *testing.T) {
		m, d := setUpAuthDriver()
		defer m.Disconnect()

		m.Mock.ExpectBegin()
		m.Mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "auths" ("created_at","updated_at","deleted_at","email","hashed_password","token") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "test@test.com", "test_hashed", sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		m.Mock.ExpectCommit()
		if err := d.Create("test@test.com", "test_hashed", ""); err != nil {
			t.Errorf("unexpected error (%v)", err)
		}
	})

	t.Run("abnormal: create auth without email", func(t *testing.T) {
		m, d := setUpAuthDriver()
		defer m.Disconnect()

		m.Mock.ExpectBegin()
		m.Mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "auths" ("created_at","updated_at","deleted_at","email","hashed_password","token") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "", "test_hashed", sqlmock.AnyArg())
		m.Mock.ExpectRollback()
		err := d.Create("", "test_hashed", "")
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
		m, d := setUpAuthDriver()
		defer m.Disconnect()

		err := d.Create("test@test.com", "", "")
		if err == nil {
			t.Errorf("expected error but got nil")
			return
		}
		var ve validator.ValidationErrors
		if !errors.As(err.Unwrap(), &ve) {
			t.Errorf("unexpected error (%v)", err)
			return
		}
	})

	t.Run("normal: fetch auth", func(t *testing.T) {
		m, d := setUpAuthDriver()
		defer m.Disconnect()

		m.Mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "auths" WHERE (email = $1 AND hashed_password = $2) AND "auths"."deleted_at" IS NULL ORDER BY "auths"."id" LIMIT $3`)).
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
		m, d := setUpAuthDriver()
		defer m.Disconnect()

		m.Mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "auths" WHERE (email = $1 AND hashed_password = $2) AND "auths"."deleted_at" IS NULL ORDER BY "auths"."id" LIMIT $3`)).
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
		m, d := setUpAuthDriver()
		defer m.Disconnect()

		m.Mock.ExpectBegin()
		m.Mock.ExpectExec(regexp.QuoteMeta(`UPDATE "auths" SET "deleted_at"=$1 WHERE (email = $2 AND hashed_password = $3) AND "auths"."deleted_at" IS NULL`)).
			WithArgs(sqlmock.AnyArg(), "test@test.com", "test_hashed").
			WillReturnResult(sqlmock.NewResult(int64(1), 1))
		m.Mock.ExpectCommit()

		if err := d.DeleteAuth("test@test.com", "test_hashed"); err != nil {
			t.Errorf("unexpected error (%v)", err)
		}
	})
}
