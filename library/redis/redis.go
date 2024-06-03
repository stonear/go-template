package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stonear/go-template/config"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	REDIS_DEFAULT_PATH = "$"
	// sets the key only if it does not already exist
	REDIS_SET_MODE_NX = "NX"
)

type Client struct {
	*redis.Client
}

func New(lc fx.Lifecycle, config *config.Config, log *otelzap.Logger) *Client {
	ctx := context.Background()

	envHost := config.RedisHost
	envPort := config.RedisPort
	envPassword := config.RedisPassword
	envDatabase := config.RedisDb

	addr := fmt.Sprintf("%s:%d", envHost, envPort)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: envPassword,
		DB:       envDatabase,
	})

	_, err := rdb.Ping(ctx).Result()
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

func (c *Client) Remember(ctx context.Context, key string, exp time.Duration, dest any, defaultFunc func() (any, error)) (any, error) {
	val, err := c.JSONGet(ctx, key, REDIS_DEFAULT_PATH).Result()
	if err == nil {
		var destArr = []any{dest}
		err = json.Unmarshal([]byte(val), &destArr)
		if len(destArr) > 0 && err == nil {
			return destArr[0], nil
		}
	}

	dest, err = defaultFunc()
	if err != nil {
		return dest, err
	}

	_, err = c.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		err = rdb.JSONSetMode(ctx, key, REDIS_DEFAULT_PATH, dest, REDIS_SET_MODE_NX).Err()
		if err != nil {
			return err
		}
		return rdb.Expire(ctx, key, exp).Err()
	})

	return dest, err
}
