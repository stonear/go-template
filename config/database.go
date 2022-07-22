package config

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func Database() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// if env database is set, use it
	envConnection := os.Getenv("DB_CONNECTION")

	if envConnection == "pgsql" {
		envHost := os.Getenv("DB_HOST")
		envPort := os.Getenv("DB_PORT")
		envDatabase := os.Getenv("DB_DATABASE")
		envUsername := os.Getenv("DB_USERNAME")
		envPassword := os.Getenv("DB_PASSWORD")

		connStr := "postgresql://" + envUsername + ":" + envPassword + "@" + envHost + ":" + envPort + "/" + envDatabase + "?sslmode=disable"

		db, err := sql.Open("postgres", connStr)
		if err != nil {
			log.Fatal(err)
		}

		return db
	} else if envConnection == "mysql" {
		// to do
		return nil
	}

	return nil
}
