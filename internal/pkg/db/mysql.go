//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest generate

package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func InitMysql(dsn string) *sql.DB {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	return db
}
