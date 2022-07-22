package database

import (
	"database/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
	//"github.com/jackc/pgx/v4"
	"github.com/stonear/go-template/helper"
	"os"
)

func New() *sql.DB {
	envHost := os.Getenv("DB_HOST")
	envPort := os.Getenv("DB_PORT")
	envDatabase := os.Getenv("DB_DATABASE")
	envUsername := os.Getenv("DB_USERNAME")
	envPassword := os.Getenv("DB_PASSWORD")

	connStr := "postgresql://" + envUsername + ":" + envPassword + "@" + envHost + ":" + envPort + "/" + envDatabase + "?sslmode=disable"
	db, err := sql.Open(os.Getenv("DB_CONNECTION"), connStr)
	helper.Panic(err)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			helper.Panic(err)
		}
	}(db)

	return db
}
