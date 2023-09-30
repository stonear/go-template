package postgres

import (
	"context"
	"os"
	"time"

	pgxZap "github.com/jackc/pgx-zap"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func New(lc fx.Lifecycle, log *zap.Logger) *pgxpool.Pool {
	ctx := context.Background()

	envHost := os.Getenv("DB_HOST")
	envPort := os.Getenv("DB_PORT")
	envDatabase := os.Getenv("DB_DATABASE")
	envUsername := os.Getenv("DB_USERNAME")
	envPassword := os.Getenv("DB_PASSWORD")

	connStr := "postgres://" + envUsername + ":" + envPassword + "@" + envHost + ":" + envPort + "/" + envDatabase + "?sslmode=disable"

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
	config.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   pgxZap.NewLogger(log),
		LogLevel: tracelog.LogLevelTrace,
	}

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
