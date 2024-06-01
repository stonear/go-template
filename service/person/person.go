package person

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

type Service interface {
	Index(c *gin.Context)
	Show(c *gin.Context)
	Store(c *gin.Context)
	Update(c *gin.Context)
	Destroy(c *gin.Context)
	Report(c *gin.Context)
}

func New(
	log *otelzap.Logger,
	pool *pgxpool.Pool,
	personDb *person.Queries,
) Service {
	return &service{
		log:      log,
		pool:     pool,
		personDb: personDb,
	}
}

type service struct {
	log      *otelzap.Logger
	pool     *pgxpool.Pool
	personDb *person.Queries
}

// @Summary		Person List
// @Description	get person list
// @Tags		person
// @Accept		json
// @Produce		json
// @Success		200		{object}	response.Response{data=[]person.IndexRow}
// @Failure		500		{object}	response.Response{data=nil}
// @Router		/person	[get]
func (s *service) Index(c *gin.Context) {
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

// @Summary		Person Show
// @Description	get person by id
// @Tags		person
// @Accept		json
// @Produce		json
// @Success		200				{object}	response.Response{data=person.ShowRow}
// @Failure		400,404,500		{object}	response.Response{data=nil}
// @Param 		id	path	string	true	"Person ID"
// @Router		/person/{id}	[get]
func (s *service) Show(c *gin.Context) {
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

// @Summary		Person Store
// @Description	save person
// @Tags		person
// @Accept		json
// @Produce		json
// @Param 		request			body		StoreRequest true "request body"
// @Success		200				{object}	response.Response{data=StoreResponse}
// @Failure		400,500			{object}	response.Response{data=nil}
// @Router		/person			[post]
func (s *service) Store(c *gin.Context) {
	ctx := c.Request.Context()

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
	var resp = StoreResponse{
		Id: person.ID.String(),
	}
	c.JSON(http.StatusOK, response.New(
		response.CodeSuccess,
		resp,
	))
}

// @Summary		Person Update
// @Description	update person
// @Tags		person
// @Accept		json
// @Produce		json
// @Param 		id 				path 		string true "User ID"
// @Param 		request			body		UpdateRequest true "request body"
// @Success		200				{object}	response.Response{data=person.Person{createdAt=time.Time,updatedAt=time.Time}}
// @Failure		400,500			{object}	response.Response{data=nil}
// @Router		/person/{id}	[put]
func (s *service) Update(c *gin.Context) {
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

// @Summary		Person Destroy
// @Description	delete person
// @Tags		person
// @Accept		json
// @Produce		json
// @Param 		id 				path 		string true "User ID"
// @Success		200				{object}	response.Response{data=nil}
// @Failure		400,500			{object}	response.Response{data=nil}
// @Router		/person/{id}	[delete]
func (s *service) Destroy(c *gin.Context) {
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

// @Summary		Person Report
// @Description	person list in pdf
// @Tags		person
// @Accept		json
// @Produce		application/pdf
// @Success		200				{file}		file
// @Failure		500				{object}	response.Response{data=nil}
// @Router		/person/report	[get]
func (s *service) Report(c *gin.Context) {
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

	c.Header("Content-Disposition", "attachment; filename=person-list.pdf")
	c.Data(http.StatusOK, "application/pdf", buffer.Bytes())
}
