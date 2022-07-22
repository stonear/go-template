package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/stonear/go-template/app/http/controllers"
)

func Api() *gin.Engine {
	router := gin.Default()

	hello := controller.Hello{}
	router.GET("/", hello.Default)
	router.GET("/hello/:name", hello.Name)

	return router
}
