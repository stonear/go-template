package main

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/stonear/go-template/controller"
	"github.com/stonear/go-template/database"
	"github.com/stonear/go-template/helper"
	"github.com/stonear/go-template/repository"
	"github.com/stonear/go-template/router"
	"github.com/stonear/go-template/service"
)

func main() {
	err := godotenv.Load()
	helper.Panic(err)

	db := database.New()
	defer func(db *sql.DB) {
		err := db.Close()
		helper.Panic(err)
	}(db)

	personRepository := repository.New()
	personService := service.New(personRepository, db)
	personController := controller.New(personService)

	r := router.New(personController)
	server := http.Server{
		Addr:    os.Getenv("APP_URL"),
		Handler: r,
	}

	err = server.ListenAndServe()
	helper.Panic(err)
}
