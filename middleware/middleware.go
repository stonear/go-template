package middleware

import (
	"aidanwoods.dev/go-paseto"
	"github.com/gin-gonic/gin"
	"github.com/stonear/go-template/config"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

type Middleware interface {
	Auth() gin.HandlerFunc
}

func New(
	config *config.Config,
	log *otelzap.Logger,
) Middleware {
	key, err := paseto.V4SymmetricKeyFromBytes([]byte(config.AppKey))
	if err != nil {
		log.Fatal("failed to parse paseto symmetric key", zap.Error(err))
	}

	return &middleware{
		config: config,
		log:    log,
		key:    key,
	}
}

type middleware struct {
	config *config.Config
	log    *otelzap.Logger

	// for auth
	key paseto.V4SymmetricKey
}
