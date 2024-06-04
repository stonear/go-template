package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/exaring/otelpgx"
	pgxZap "github.com/jackc/pgx-zap"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/stonear/go-template/config"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func New(lc fx.Lifecycle, config *config.Config, log *otelzap.Logger) *pgxpool.Pool {
	ctx := context.Background()

	envHost := config.DbHost
	envPort := config.DbPort
	envDatabase := config.DbDatabase
	envUsername := config.DbUsername
	envPassword := config.DbPassword

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", envUsername, envPassword, envHost, envPort, envDatabase)

	pgxConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatal("[Database] Unable to parse connection url", zap.Error(err))
	}
	pgxConfig.MaxConns = 50
	pgxConfig.MinConns = 0
	pgxConfig.MaxConnLifetime = time.Hour
	pgxConfig.MaxConnIdleTime = time.Minute * 30
	pgxConfig.HealthCheckPeriod = time.Minute
	pgxConfig.ConnConfig.ConnectTimeout = time.Second * 10
	pgxConfig.ConnConfig.RuntimeParams["timezone"] = config.DbTz
	pgxConfig.ConnConfig.Tracer = NewTracer(
		// zap
		&tracelog.TraceLog{
			Logger:   pgxZap.NewLogger(log.Logger),
			LogLevel: tracelog.LogLevelTrace,
		},
		// opentelemetry
		otelpgx.NewTracer(),
	)

	conn, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		log.Fatal("[Database] Unable to connect to database", zap.Error(err))
	}

	err = conn.Ping(ctx)
	if err != nil {
		log.Fatal("[Database] Unable to ping database", zap.Error(err))
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			conn.Close()
			log.Info("[Database] DB connection is stopped")
			return nil
		},
	})

	return conn
}
