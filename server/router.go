package server

import (
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
	"github.com/stonear/go-template/docs"
	"github.com/stonear/go-template/service/person"
	"github.com/stonear/go-template/service/pokemon"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router(
	r *gin.Engine,
	personSvc person.Service,
	pokemonSvc pokemon.Service,
) {
	basePath := "/v1"
	docs.SwaggerInfo.BasePath = basePath

	securityConfig := secure.DefaultConfig()
	securityConfig.SSLRedirect = false

	api := r.Group(basePath).Use(secure.New(securityConfig))
	{
		api.GET("/person", personSvc.Index)
		api.GET("/person/report", personSvc.Report)
		api.GET("/person/:id", personSvc.Show)
		api.POST("/person", personSvc.Store)
		api.PUT("/person/:id", personSvc.Update)
		api.DELETE("/person/:id", personSvc.Destroy)

		api.GET("/pokemon", pokemonSvc.Index)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
