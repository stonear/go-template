package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stonear/go-template/db/person"
	"github.com/stonear/go-template/response"
	"github.com/stonear/go-template/validator"
	"go.uber.org/zap"
)

type PersonService interface {
	Index(ctx *gin.Context)
	Show(ctx *gin.Context)
	Store(ctx *gin.Context)
	Update(ctx *gin.Context)
	Destroy(ctx *gin.Context)
}

func NewPersonService(
	log *zap.Logger,
	pool *pgxpool.Pool,
	personDb *person.Queries,
) PersonService {
	return &personService{
		log:      log,
		pool:     pool,
		personDb: personDb,
	}
}

type personService struct {
	log      *zap.Logger
	pool     *pgxpool.Pool
	personDb *person.Queries
}

func (s *personService) Index(ctx *gin.Context) {
	persons, err := s.personDb.Index(ctx, s.pool)
	if err != nil {
		s.log.Error("failed to get personDb.Index", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, response.New(
			response.CodeGeneralError,
			nil,
		))
		return
	}
	ctx.JSON(http.StatusOK, response.New(
		response.CodeSuccess,
		persons,
	))
}

func (s *personService) Show(ctx *gin.Context) {
	type ShowRequest struct {
		Id int `uri:"id" binding:"required,gt=0"`
	}
	var req ShowRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.log.Error("failed to bind request", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, response.New(
			response.CodeInvalidFormat,
			validator.Message(err),
		))
		return
	}
	person, err := s.personDb.Show(ctx, s.pool, int64(req.Id))
	if err != nil {
		s.log.Error(
			"failed to get personDb.Show",
			zap.Int("Id", req.Id),
			zap.Error(err),
		)
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, response.New(
				response.CodeNotFound,
				nil,
			))
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.New(
			response.CodeGeneralError,
			nil,
		))
		return
	}
	ctx.JSON(http.StatusOK, response.New(
		response.CodeSuccess,
		person,
	))
}

func (s *personService) Store(ctx *gin.Context) {
	type StoreRequest struct {
		Name string `form:"name" binding:"required,notEvil"`
	}
	var req StoreRequest
	err := ctx.ShouldBindWith(&req, binding.JSON)
	if err != nil {
		s.log.Error("failed to bind request", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, response.New(
			response.CodeInvalidFormat,
			validator.Message(err),
		))
		return
	}
	person, err := s.personDb.Store(ctx, s.pool, req.Name)
	if err != nil {
		s.log.Error(
			"failed to call personDb.Store",
			zap.String("Name", req.Name),
			zap.Error(err),
		)
		ctx.JSON(http.StatusInternalServerError, response.New(
			response.CodeGeneralError,
			nil,
		))
		return
	}
	type StoreResponse struct {
		Id int `json:"id"`
	}
	var resp = StoreResponse{
		Id: int(person.ID),
	}
	ctx.JSON(http.StatusOK, response.New(
		response.CodeSuccess,
		resp,
	))
}

func (s *personService) Update(ctx *gin.Context) {
	// TODO: implement this
	ctx.JSON(http.StatusOK, response.New(
		response.CodeSuccess,
		nil,
	))
}

func (s *personService) Destroy(ctx *gin.Context) {
	// TODO: implement this
	ctx.JSON(http.StatusOK, response.New(
		response.CodeSuccess,
		nil,
	))
}
