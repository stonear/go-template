package controller

import (
	"net/http"
	"strconv"

	"github.com/stonear/go-template/helper"

	"github.com/gin-gonic/gin"
	"github.com/stonear/go-template/entity"
	"github.com/stonear/go-template/service"
)

type Controller interface {
	Index(ctx *gin.Context)
	Show(ctx *gin.Context)
	Store(ctx *gin.Context)
	Update(ctx *gin.Context)
	Destroy(ctx *gin.Context)
}

func New(serv service.Service) Controller {
	return &controller{
		Service: serv,
	}
}

type controller struct {
	Service service.Service
}

func (c controller) Index(ctx *gin.Context) {
	persons := c.Service.Index(ctx)
	ctx.JSON(http.StatusOK, gin.H{"persons": persons})
}

func (c controller) Show(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	helper.Panic(err)
	person := c.Service.Show(ctx, id)
	ctx.JSON(http.StatusOK, gin.H{"person": person})
}

func (c *controller) Store(ctx *gin.Context) {
	person := entity.Person{}
	err := ctx.ShouldBind(&person)
	helper.Panic(err)
	id, err := c.Service.Store(ctx, person)
	helper.Panic(err)
	ctx.JSON(http.StatusCreated, gin.H{"id": id})
}

func (c *controller) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	helper.Panic(err)
	person := entity.Person{}
	err = ctx.ShouldBind(&person)
	helper.Panic(err)
	person, err = c.Service.Update(ctx, id, person)
	helper.Panic(err)
	ctx.JSON(http.StatusOK, gin.H{"person": person})
}

func (c *controller) Destroy(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	helper.Panic(err)
	err = c.Service.Destroy(ctx, id)
	helper.Panic(err)
	ctx.Writer.WriteHeader(http.StatusOK)
}
