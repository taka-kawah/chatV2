package db

import (
	"database/sql"
	"fmt"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type mockDbInstances struct {
	GormDb *gorm.DB
	SqlDb  *sql.DB
	Mock   sqlmock.Sqlmock
}

func NewMockDbInstances() (*mockDbInstances, error) {
	sqlDb, mock, err := sqlmock.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create mock db (%s)", err)
	}
	gormDb, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDb}), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create gorm mock (%s)", err)
	}
	return &mockDbInstances{GormDb: gormDb, SqlDb: sqlDb, Mock: mock}, nil
}

func (db *mockDbInstances) Disconnect() error {
	if err := db.SqlDb.Close(); err != nil {
		return fmt.Errorf("failed to disconnect mock db")
	}
	return nil
}
