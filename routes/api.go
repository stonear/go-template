package routes

import (
	"github.com/julienschmidt/httprouter"
	controller "github.com/stonear/go-template/app/http/controllers"
)

func Api() *httprouter.Router {
	router := httprouter.New()

	hello := controller.Hello{}
	router.GET("/", hello.Default)
	router.GET("/hello/:name", hello.Name)

	return router
}
