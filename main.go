package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stonear/go-template/client/pokemon"
	"github.com/stonear/go-template/config"
	"github.com/stonear/go-template/db/person"
	"github.com/stonear/go-template/library/httpclient"
	"github.com/stonear/go-template/library/postgres"
	"github.com/stonear/go-template/library/redis"
	"github.com/stonear/go-template/logger"
	"github.com/stonear/go-template/server"
	personService "github.com/stonear/go-template/service/person"
	pokemonService "github.com/stonear/go-template/service/pokemon"
	"github.com/stonear/go-template/tracer"
	"github.com/stonear/go-template/validator"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

// @title						Go-Template
// @version						1.0.0
// @description					This is an example.
// @termsOfService				http://swagger.io/terms/
// @host						localhost:8080
// @BasePath					/v1
// @securityDefinitions.apikey	Bearer
// @in							header
// @name						Authorization
// @description					Type "Bearer" followed by a space and JWT token.
func main() {
	fx.New(
		fx.WithLogger(func(log *otelzap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{
				Logger: log.Logger,
			}
		}),
		fx.Provide(
			logger.New,
			postgres.New,
			redis.New,
			httpclient.New,

			person.New,
			pokemon.New,

			personService.New,
			pokemonService.New,

			server.New,
		),
		fx.Invoke(
			config.Load,
			validator.Load,
			tracer.Load,
			func(*pgxpool.Pool) {},
			server.Router,
			func(*gin.Engine) {},
		),
	).Run()
}
