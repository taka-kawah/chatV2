package db

import (
	"back/interfaces"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DbInstances struct {
	GormDb *gorm.DB
	SqlDB  *sql.DB
}

func NewDbInstances() (*DbInstances, interfaces.CustomError) {
	dsn, err := loadDsn()
	if err != nil {
		return nil, err.(interfaces.CustomError)
	}

	gormDb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{PrepareStmt: true})
	if err != nil {
		return nil, &ConnectionToDbError{msg: "failed to open gormDb", err: err}
	}

	sqlDb, err := gormDb.DB()
	if err != nil {
		return nil, &ConnectionToDbError{msg: "failed to get SQL instance", err: err}
	}
	return &DbInstances{GormDb: gormDb, SqlDB: sqlDb}, nil
}

func (db *DbInstances) Disconnect() error {
	if err := db.SqlDB.Close(); err != nil {
		return &ConnectionToDbError{msg: "failed to close sqlDb", err: err}
	}
	return nil
}

func loadDsn() (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", &ConnectionToDbError{msg: "failed to load dotenv", err: err}
	}

	env := map[string]string{
		"host":     os.Getenv("DB_HOST"),
		"user":     os.Getenv("DB_USER"),
		"password": os.Getenv("DB_PASSWORD"),
		"dbname":   os.Getenv("DB_NAME"),
		"port":     os.Getenv("DB_PORT"),
	}
	emptyEnvs := make([]string, 0, 5)
	for key, val := range env {
		if len(val) == 0 {
			emptyEnvs = append(emptyEnvs, key)
		}
	}
	if len(emptyEnvs) > 0 {
		return "", &ConnectionToDbError{msg: "missing required environment variables", err: errors.New(strings.Join(emptyEnvs, ","))}
	}

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=Asia/Tokyo sslmode=disable", env["host"], env["user"], env["password"], env["dbname"], env["port"]), nil
}

type ConnectionToDbError struct {
	msg string
	err error
}

func (e *ConnectionToDbError) Error() string {
	return fmt.Sprintf("error in establishing connection to db %s (%s)", e.msg, e.err)
}

func (e *ConnectionToDbError) Unwrap() error {
	return e.err
}
