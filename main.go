package main

import (
	"github.com/stonear/go-template/routes"
)

func main() {
	router := routes.Api()
	router.Run(":8080")
}
