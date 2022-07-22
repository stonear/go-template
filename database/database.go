package database

import (
	"database/sql"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stonear/go-template/helper"
)

func New() *sql.DB {
	envHost := os.Getenv("DB_HOST")
	envPort := os.Getenv("DB_PORT")
	envDatabase := os.Getenv("DB_DATABASE")
	envUsername := os.Getenv("DB_USERNAME")
	envPassword := os.Getenv("DB_PASSWORD")

	driver := os.Getenv("DB_CONNECTION")
	connStr := ""
	if driver == "pgx" {
		connStr = "postgresql://" + envUsername + ":" + envPassword + "@" + envHost + ":" + envPort + "/" + envDatabase + "?sslmode=disable"
	}

	// TODO: add support for other drivers

	db, err := sql.Open(driver, connStr)
	helper.Panic(err)

	return db
}
