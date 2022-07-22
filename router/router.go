package router

import (
	"github.com/gin-gonic/gin"
	"github.com/stonear/go-template/controller"
)

func New(c controller.Controller) *gin.Engine {
	r := gin.Default()

	r.GET("/api/person", c.Index)
	r.GET("/api/person/:id", c.Show)

	return r
}
