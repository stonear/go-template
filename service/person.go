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
	type ShowUri struct {
		Id int `uri:"id" binding:"required,gt=0"`
	}
	var uri ShowUri
	err := ctx.ShouldBindUri(&uri)
	if err != nil {
		s.log.Error("failed to bind request", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, response.New(
			response.CodeInvalidFormat,
			validator.Message(err),
		))
		return
	}
	person, err := s.personDb.Show(ctx, s.pool, int64(uri.Id))
	if err != nil {
		s.log.Error(
			"failed to get personDb.Show",
			zap.Int("Id", uri.Id),
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
	type UpdateUri struct {
		Id int `uri:"id" binding:"required,gt=0"`
	}
	var uri UpdateUri
	err := ctx.ShouldBindUri(&uri)
	if err != nil {
		s.log.Error("failed to bind request", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, response.New(
			response.CodeInvalidFormat,
			validator.Message(err),
		))
		return
	}
	type UpdateRequest struct {
		Name string `form:"name" binding:"required,notEvil"`
	}
	var req UpdateRequest
	err = ctx.ShouldBindWith(&req, binding.JSON)
	if err != nil {
		s.log.Error("failed to bind request", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, response.New(
			response.CodeInvalidFormat,
			validator.Message(err),
		))
		return
	}
	person, err := s.personDb.Update(ctx, s.pool, person.UpdateParams{
		ID:   int64(uri.Id),
		Name: req.Name,
	})
	if err != nil {
		s.log.Error(
			"failed to call personDb.Update",
			zap.Int("Id", uri.Id),
			zap.String("Name", req.Name),
			zap.Error(err),
		)
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

func (s *personService) Destroy(ctx *gin.Context) {
	type DestroyUri struct {
		Id int `uri:"id" binding:"required,gt=0"`
	}
	var uri DestroyUri
	err := ctx.ShouldBindUri(&uri)
	if err != nil {
		s.log.Error("failed to bind request", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, response.New(
			response.CodeInvalidFormat,
			validator.Message(err),
		))
		return
	}
	err = s.personDb.Destroy(ctx, s.pool, int64(uri.Id))
	if err != nil {
		s.log.Error(
			"failed to call personDb.Destroy",
			zap.Int("Id", uri.Id),
			zap.Error(err),
		)
		ctx.JSON(http.StatusInternalServerError, response.New(
			response.CodeGeneralError,
			nil,
		))
		return
	}
	ctx.JSON(http.StatusOK, response.New(
		response.CodeSuccess,
		nil,
	))
}
