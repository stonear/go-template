package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Hello struct {
	// to do add properties
}

func (h *Hello) Default(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello world!\n",
	})
}

func (h *Hello) Name(c *gin.Context) {
	name := c.Param("name")
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello " + name + "!\n",
	})
}
