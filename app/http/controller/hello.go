package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HelloController struct {
	// to do add properties
}

func (h *HelloController) Default(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello world!\n",
	})
}

func (h *HelloController) Name(c *gin.Context) {
	name := c.Param("name")
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello " + name + "!\n",
	})
}
