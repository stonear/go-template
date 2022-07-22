package controller

import (
	"github.com/stonear/go-template/helper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/stonear/go-template/entity"
	"github.com/stonear/go-template/service"
)

type Controller interface {
	Create(ctx *gin.Context, person entity.Person)
	Update(ctx *gin.Context, person entity.Person)
	Delete(ctx *gin.Context, id int)
	FindById(ctx *gin.Context)
	FindAll(ctx *gin.Context)
}

func New(serv service.Service) Controller {
	return &controller{
		Service: serv,
	}
}

type controller struct {
	Service service.Service
}

func (c *controller) Create(ctx *gin.Context, person entity.Person) {
	//TODO implement me
	panic("implement me")
}

func (c *controller) Update(ctx *gin.Context, person entity.Person) {
	//TODO implement me
	panic("implement me")
}

func (c *controller) Delete(ctx *gin.Context, id int) {
	//TODO implement me
	panic("implement me")
}

func (c controller) FindById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	helper.Panic(err)
	person := c.Service.FindById(ctx, id)
	ctx.JSON(http.StatusOK, gin.H{"person": person})
}

func (c controller) FindAll(ctx *gin.Context) {
	persons := c.Service.FindAll(ctx)
	ctx.JSON(http.StatusOK, gin.H{"persons": persons})
}
