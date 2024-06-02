package auth

import (
	"net/http"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stonear/go-template/config"
	"github.com/stonear/go-template/db/auth"
	"github.com/stonear/go-template/response"
	"github.com/stonear/go-template/validator"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Login(c *gin.Context)
}

func New(
	config *config.Config,
	log *otelzap.Logger,
	pool *pgxpool.Pool,
	authDb *auth.Queries,
) Service {
	return &service{
		config: config,
		log:    log,
		pool:   pool,
		authDb: authDb,
	}
}

type service struct {
	config *config.Config
	log    *otelzap.Logger
	pool   *pgxpool.Pool
	authDb *auth.Queries
}

// @Summary		Auth Login
// @Description	login
// @Tags		auth
// @Accept		json
// @Produce		json
// @Param 		request			body		LoginRequest true "request body"
// @Success		200				{object}	response.Response{data=string}
// @Failure		403,500			{object}	response.Response{data=nil}
// @Router		/auth/login		[post]
func (s *service) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var req LoginRequest
	err := c.ShouldBindWith(&req, binding.JSON)
	if err != nil {
		s.log.Ctx(ctx).Error("failed to bind request", zap.Error(err))
		c.JSON(http.StatusBadRequest, response.New(
			response.CodeInvalidFormat,
			validator.Message(err),
		))
		return
	}

	user, err := s.authDb.Show(ctx, s.pool, req.Username)
	if err != nil {
		s.log.Ctx(ctx).Error(
			"failed to get authDb.Show",
			zap.String("username", req.Username),
			zap.Error(err),
		)
		c.JSON(http.StatusForbidden, response.New(
			response.CodeNotPermitted,
			nil,
		))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		s.log.Ctx(ctx).Error(
			"failed to compare bcrypt hash and password",
			zap.String("hash", user.Password),
			zap.Error(err),
		)
		c.JSON(http.StatusForbidden, response.New(
			response.CodeNotPermitted,
			nil,
		))
		return
	}

	token := paseto.NewToken()
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(2 * time.Hour))
	token.SetString("id", user.ID.String())
	token.SetString("username", user.Username)

	key, err := paseto.V4SymmetricKeyFromBytes([]byte(s.config.AppKey))
	if err != nil {
		s.log.Ctx(ctx).Error("failed to parse paseto symmetric key", zap.Error(err))
		c.JSON(http.StatusBadRequest, response.New(
			response.CodeGeneralError,
			validator.Message(err),
		))
		return
	}

	encrypted := token.V4Encrypt(key, nil)

	c.JSON(http.StatusOK, response.New(
		response.CodeSuccess,
		encrypted,
	))
}
