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
	Create(ctx *gin.Context, person entity.Person)
	Update(ctx *gin.Context, person entity.Person)
	Destroy(ctx *gin.Context, id int)
}

func New(serv service.Service) Controller {
	return &controller{
		Service: serv,
	}
}

type controller struct {
	Service service.Service
}

func (c controller) Show(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	helper.Panic(err)
	person := c.Service.Show(ctx, id)
	ctx.JSON(http.StatusOK, gin.H{"person": person})
}

func (c controller) Index(ctx *gin.Context) {
	persons := c.Service.Index(ctx)
	ctx.JSON(http.StatusOK, gin.H{"persons": persons})
}

func (c *controller) Create(ctx *gin.Context, person entity.Person) {
	//TODO implement me
	panic("implement me")
}

func (c *controller) Update(ctx *gin.Context, person entity.Person) {
	//TODO implement me
	panic("implement me")
}

func (c *controller) Destroy(ctx *gin.Context, id int) {
	//TODO implement me
	panic("implement me")
}
