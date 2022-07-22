package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/stonear/go-template/app/http/controller"
	"github.com/stonear/go-template/config"
)

func Api() *gin.Engine {
	router := gin.Default()

	hello := controller.HelloController{}
	router.GET("/", hello.Default)
	router.GET("/hello/:name", hello.Name)

	mhs := controller.MhsController{
		Db: config.Database(),
	}
	router.GET("/mhs", mhs.Index)

	return router
}
