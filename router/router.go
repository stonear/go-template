package router

import (
	"github.com/gin-gonic/gin"
	"github.com/stonear/go-template/controller"
)

func New(c controller.Controller) *gin.Engine {
	r := gin.Default()

	api := r.Group("/v1")
	{
		api.GET("/person", c.Index)
		api.GET("/person/:id", c.Show)
		api.POST("/person", c.Store)
	}

	return r
}
