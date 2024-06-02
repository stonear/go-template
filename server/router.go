package server

import (
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
	"github.com/stonear/go-template/docs"
	"github.com/stonear/go-template/middleware"
	"github.com/stonear/go-template/service/auth"
	"github.com/stonear/go-template/service/person"
	"github.com/stonear/go-template/service/pokemon"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router(
	r *gin.Engine,
	middleware middleware.Middleware,
	authSvc auth.Service,
	personSvc person.Service,
	pokemonSvc pokemon.Service,
) {
	basePath := "/v1"
	docs.SwaggerInfo.BasePath = basePath

	securityConfig := secure.DefaultConfig()
	securityConfig.SSLRedirect = false

	api := r.Group(basePath)
	api.Use(secure.New(securityConfig))
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", authSvc.Login)
		}

		protected := api
		protected.Use(middleware.Auth())
		{
			person := protected.Group("/person")
			{
				person.GET("", personSvc.Index)
				person.GET("/report", personSvc.Report)
				person.GET("/:id", personSvc.Show)
				person.POST("", personSvc.Store)
				person.PUT("/:id", personSvc.Update)
				person.DELETE("/:id", personSvc.Destroy)
			}

			protected.GET("/pokemon", pokemonSvc.Index)
		}
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
