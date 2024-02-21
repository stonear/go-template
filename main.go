package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stonear/go-template/client/pokemon"
	"github.com/stonear/go-template/config"
	"github.com/stonear/go-template/db/person"
	"github.com/stonear/go-template/library/httpclient"
	"github.com/stonear/go-template/library/postgres"
	"github.com/stonear/go-template/logger"
	"github.com/stonear/go-template/server"
	"github.com/stonear/go-template/service"
	"github.com/stonear/go-template/tracer"
	"github.com/stonear/go-template/validator"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

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
			httpclient.New,

			person.New,
			pokemon.New,

			service.NewPersonService,
			service.NewPokemonService,

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
