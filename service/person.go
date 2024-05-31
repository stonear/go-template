package service

import (
	"bytes"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/signintech/gopdf"
	"github.com/signintech/pdft"
	"github.com/stonear/go-template/db/person"
	"github.com/stonear/go-template/response"
	"github.com/stonear/go-template/validator"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

type PersonService interface {
	Index(c *gin.Context)
	Show(c *gin.Context)
	Store(c *gin.Context)
	Update(c *gin.Context)
	Destroy(c *gin.Context)
	Report(c *gin.Context)
}

func NewPersonService(
	log *otelzap.Logger,
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
	log      *otelzap.Logger
	pool     *pgxpool.Pool
	personDb *person.Queries
}

func (s *personService) Index(c *gin.Context) {
	ctx := c.Request.Context()
	persons, err := s.personDb.Index(ctx, s.pool)
	if err != nil {
		s.log.Ctx(ctx).Error("failed to get personDb.Index", zap.Error(err))
		c.JSON(http.StatusInternalServerError, response.New(
			response.CodeGeneralError,
			nil,
		))
		return
	}
	c.JSON(http.StatusOK, response.New(
		response.CodeSuccess,
		persons,
	))
}

func (s *personService) Show(c *gin.Context) {
	ctx := c.Request.Context()
	type ShowUri struct {
		Id string `uri:"id" binding:"required,gt=0"`
	}
	var uri ShowUri
	err := c.ShouldBindUri(&uri)
	if err != nil {
		s.log.Ctx(ctx).Error("failed to bind request", zap.Error(err))
		c.JSON(http.StatusBadRequest, response.New(
			response.CodeInvalidFormat,
			validator.Message(err),
		))
		return
	}
	id, err := uuid.Parse(uri.Id)
	if err != nil {
		s.log.Ctx(ctx).Error("failed to parse uuid", zap.Error(err))
		c.JSON(http.StatusBadRequest, response.New(
			response.CodeInvalidFormat,
			validator.Message(err),
		))
		return
	}
	person, err := s.personDb.Show(ctx, s.pool, id)
	if err != nil {
		s.log.Ctx(ctx).Error(
			"failed to get personDb.Show",
			zap.String("Id", uri.Id),
			zap.Error(err),
		)
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, response.New(
				response.CodeNotFound,
				nil,
			))
			return
		}
		c.JSON(http.StatusInternalServerError, response.New(
			response.CodeGeneralError,
			nil,
		))
		return
	}
	c.JSON(http.StatusOK, response.New(
		response.CodeSuccess,
		person,
	))
}

func (s *personService) Store(c *gin.Context) {
	ctx := c.Request.Context()
	type StoreRequest struct {
		Name string `form:"name" binding:"required,notEvil"`
	}
	var req StoreRequest
	err := c.ShouldBindWith(&req, binding.JSON)
	if err != nil {
		s.log.Ctx(ctx).Error("failed to bind request", zap.Error(err))
		c.JSON(http.StatusBadRequest, response.New(
			response.CodeInvalidFormat,
			validator.Message(err),
		))
		return
	}
	person, err := s.personDb.Store(ctx, s.pool, req.Name)
	if err != nil {
		s.log.Ctx(ctx).Error(
			"failed to call personDb.Store",
			zap.String("Name", req.Name),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, response.New(
			response.CodeGeneralError,
			nil,
		))
		return
	}
	type StoreResponse struct {
		Id string `json:"id"`
	}
	var resp = StoreResponse{
		Id: person.ID.String(),
	}
	c.JSON(http.StatusOK, response.New(
		response.CodeSuccess,
		resp,
	))
}

func (s *personService) Update(c *gin.Context) {
	ctx := c.Request.Context()
	type UpdateUri struct {
		Id string `uri:"id" binding:"required,gt=0"`
	}
	var uri UpdateUri
	err := c.ShouldBindUri(&uri)
	if err != nil {
		s.log.Ctx(ctx).Error("failed to bind request", zap.Error(err))
		c.JSON(http.StatusBadRequest, response.New(
			response.CodeInvalidFormat,
			validator.Message(err),
		))
		return
	}
	type UpdateRequest struct {
		Name string `form:"name" binding:"required,notEvil"`
	}
	var req UpdateRequest
	err = c.ShouldBindWith(&req, binding.JSON)
	if err != nil {
		s.log.Ctx(ctx).Error("failed to bind request", zap.Error(err))
		c.JSON(http.StatusBadRequest, response.New(
			response.CodeInvalidFormat,
			validator.Message(err),
		))
		return
	}
	id, err := uuid.Parse(uri.Id)
	if err != nil {
		s.log.Ctx(ctx).Error("failed to parse uuid", zap.Error(err))
		c.JSON(http.StatusBadRequest, response.New(
			response.CodeInvalidFormat,
			validator.Message(err),
		))
		return
	}
	person, err := s.personDb.Update(ctx, s.pool, person.UpdateParams{
		ID:   id,
		Name: req.Name,
	})
	if err != nil {
		s.log.Ctx(ctx).Error(
			"failed to call personDb.Update",
			zap.String("Id", uri.Id),
			zap.String("Name", req.Name),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, response.New(
			response.CodeGeneralError,
			nil,
		))
		return
	}
	c.JSON(http.StatusOK, response.New(
		response.CodeSuccess,
		person,
	))
}

func (s *personService) Destroy(c *gin.Context) {
	ctx := c.Request.Context()
	type DestroyUri struct {
		Id string `uri:"id" binding:"required,gt=0"`
	}
	var uri DestroyUri
	err := c.ShouldBindUri(&uri)
	if err != nil {
		s.log.Ctx(ctx).Error("failed to bind request", zap.Error(err))
		c.JSON(http.StatusBadRequest, response.New(
			response.CodeInvalidFormat,
			validator.Message(err),
		))
		return
	}
	id, err := uuid.Parse(uri.Id)
	if err != nil {
		s.log.Ctx(ctx).Error("failed to parse uuid", zap.Error(err))
		c.JSON(http.StatusBadRequest, response.New(
			response.CodeInvalidFormat,
			validator.Message(err),
		))
		return
	}
	err = s.personDb.Destroy(ctx, s.pool, id)
	if err != nil {
		s.log.Ctx(ctx).Error(
			"failed to call personDb.Destroy",
			zap.String("Id", uri.Id),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, response.New(
			response.CodeGeneralError,
			nil,
		))
		return
	}
	c.JSON(http.StatusOK, response.New(
		response.CodeSuccess,
		nil,
	))
}

func (s *personService) Report(c *gin.Context) {
	ctx := c.Request.Context()
	persons, err := s.personDb.Index(ctx, s.pool)
	if err != nil {
		s.log.Ctx(ctx).Error("failed to get personDb.Index", zap.Error(err))
		c.JSON(http.StatusInternalServerError, response.New(
			response.CodeGeneralError,
			nil,
		))
		return
	}

	var pt pdft.PDFt
	err = pt.Open("template/blank.pdf")
	if err != nil {
		s.log.Ctx(ctx).Error("failed to load pdf template", zap.Error(err))
		c.JSON(http.StatusInternalServerError, response.New(
			response.CodeGeneralError,
			nil,
		))
		return
	}

	err = pt.AddFont("roboto", "template/Roboto-Regular.ttf")
	if err != nil {
		s.log.Ctx(ctx).Error("failed to load font", zap.Error(err))
		c.JSON(http.StatusInternalServerError, response.New(
			response.CodeGeneralError,
			nil,
		))
		return
	}

	err = pt.SetFont("roboto", "", 14)
	if err != nil {
		s.log.Ctx(ctx).Error("failed to set font", zap.Error(err))
		c.JSON(http.StatusInternalServerError, response.New(
			response.CodeGeneralError,
			nil,
		))
		return
	}

	x := float64(10)
	y := float64(10)
	for _, person := range persons {
		err = pt.Insert(person.ID.String(), 1, x, y, 100, 14, gopdf.Left|gopdf.Bottom, nil)
		if err != nil {
			s.log.Ctx(ctx).Error("failed to insert text", zap.Error(err))
			c.JSON(http.StatusInternalServerError, response.New(
				response.CodeGeneralError,
				nil,
			))
			return
		}

		err = pt.Insert(person.Name, 1, x+360, y, 100, 14, gopdf.Left|gopdf.Bottom, nil)
		if err != nil {
			s.log.Ctx(ctx).Error("failed to insert text", zap.Error(err))
			c.JSON(http.StatusInternalServerError, response.New(
				response.CodeGeneralError,
				nil,
			))
			return
		}

		y += 18
	}

	var buffer bytes.Buffer
	err = pt.SaveTo(&buffer)
	if err != nil {
		s.log.Ctx(ctx).Error("failed to save pdf", zap.Error(err))
		c.JSON(http.StatusInternalServerError, response.New(
			response.CodeGeneralError,
			nil,
		))
		return
	}

	c.Data(http.StatusOK, "application/pdf", buffer.Bytes())
}
