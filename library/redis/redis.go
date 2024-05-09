package redis

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"github.com/vmihailenco/msgpack/v5"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Client struct {
	*redis.Client
}

func New(lc fx.Lifecycle, log *otelzap.Logger) *Client {
	ctx := context.Background()

	envHost := os.Getenv("REDIS_HOST")
	envPort := os.Getenv("REDIS_PORT")
	envPassword := os.Getenv("REDIS_PASSWORD")
	envDatabase := os.Getenv("REDIS_DB")

	addr := fmt.Sprintf("%s:%s", envHost, envPort)
	DB, err := strconv.Atoi(envDatabase)
	if err != nil {
		log.Fatal("[Redis] Invalid env REDIS_DB number", zap.Error(err))
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: envPassword,
		DB:       DB,
	})

	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal("[Redis] Connection was refused", zap.Error(err))
	}

	if err := redisotel.InstrumentTracing(rdb); err != nil {
		log.Fatal("[Redis] Unable add InstrumentTracing", zap.Error(err))
	}
	if err := redisotel.InstrumentMetrics(rdb); err != nil {
		log.Fatal("[Redis] Unable add InstrumentMetrics", zap.Error(err))
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			err := rdb.Close()
			log.Info("[Redis] Connection is stopped")
			return err
		},
	})

	return &Client{rdb}
}

func (c *Client) Remember(ctx context.Context, key string, exp time.Duration, destType any, defaultFunc func() (any, error)) (any, error) {
	dest := destType

	val, err := c.Get(ctx, key).Result()
	if err == nil {
		err = msgpack.Unmarshal([]byte(val), dest)
		return dest, err
	}

	dest, err = defaultFunc()
	if err != nil {
		return dest, err
	}

	bytes, err := msgpack.Marshal(dest)
	if err != nil {
		return dest, err
	}

	return dest, c.Set(ctx, key, bytes, exp).Err()
}
