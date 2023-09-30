package server

import (
	"github.com/gin-gonic/gin"
	"github.com/stonear/go-template/service"
)

func Router(
	r *gin.Engine,
	personSvc service.PersonService,
) {
	api := r.Group("/v1")
	{
		api.GET("/person", personSvc.Index)
		api.GET("/person/:id", personSvc.Show)
		api.POST("/person", personSvc.Store)
		api.PUT("/person/:id", personSvc.Update)
		api.DELETE("/person/:id", personSvc.Destroy)
	}
}
