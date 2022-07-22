package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/stonear/go-template/routes"
)

func main() {
	router := routes.Api()

	err := http.ListenAndServe(":8080", router)
	if errors.Is(err, http.ErrServerClosed) {
		log.Print("Server closed\n")
	} else {
		log.Fatalf("Error: %v\n", err)
	}
}
