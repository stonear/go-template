package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stonear/go-template/config"
	"github.com/stonear/go-template/db/person"
	"github.com/stonear/go-template/library/postgres"
	"github.com/stonear/go-template/logger"
	"github.com/stonear/go-template/server"
	"github.com/stonear/go-template/service"
	"github.com/stonear/go-template/validator"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{
				Logger: log,
			}
		}),
		fx.Provide(
			logger.New,
			postgres.New,
			person.New,
			service.NewPersonService,
			server.New,
		),
		fx.Invoke(
			config.Load,
			validator.Load,
			func(*pgxpool.Pool) {},
			server.Router,
			func(*gin.Engine) {},
		),
	).Run()
}
