package postgres

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/exaring/otelpgx"
	pgxZap "github.com/jackc/pgx-zap"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func New(lc fx.Lifecycle, log *otelzap.Logger) *pgxpool.Pool {
	ctx := context.Background()

	envHost := os.Getenv("DB_HOST")
	envPort := os.Getenv("DB_PORT")
	envDatabase := os.Getenv("DB_DATABASE")
	envUsername := os.Getenv("DB_USERNAME")
	envPassword := os.Getenv("DB_PASSWORD")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", envUsername, envPassword, envHost, envPort, envDatabase)

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatal("[Database] Unable to parse connection url", zap.Error(err))
	}
	config.MaxConns = 50
	config.MinConns = 0
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = time.Minute * 30
	config.HealthCheckPeriod = time.Minute
	config.ConnConfig.ConnectTimeout = time.Second * 10
	config.ConnConfig.RuntimeParams["timezone"] = os.Getenv("DB_TZ")
	config.ConnConfig.Tracer = NewTracer(
		// zap
		&tracelog.TraceLog{
			Logger:   pgxZap.NewLogger(log.Logger),
			LogLevel: tracelog.LogLevelTrace,
		},
		// opentelemetry
		otelpgx.NewTracer(),
	)

	conn, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatal("[Database] Unable to connect to database", zap.Error(err))
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
