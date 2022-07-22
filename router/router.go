package router

import (
	"github.com/gin-gonic/gin"
	"github.com/stonear/go-template/controller"
)

func New(c controller.Controller) *gin.Engine {
	r := gin.Default()

	r.GET("/api/persons", c.FindAll)

	return r
}
