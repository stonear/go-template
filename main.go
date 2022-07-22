package main

import (
	"github.com/joho/godotenv"
	"github.com/stonear/go-template/controller"
	"github.com/stonear/go-template/database"
	"github.com/stonear/go-template/helper"
	"github.com/stonear/go-template/repository"
	"github.com/stonear/go-template/router"
	"github.com/stonear/go-template/service"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load()
	helper.Panic(err)

	db := database.New()
	personRepository := repository.New()
	personService := service.New(personRepository, db)
	personController := controller.New(personService)

	r := router.New(personController)
	server := http.Server{
		Addr:    os.Getenv("APP_URL"),
		Handler: r,
	}

	if err := server.ListenAndServe(); err != nil {
		helper.Panic(err)
	}
}
