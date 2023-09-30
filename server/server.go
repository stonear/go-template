package server

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/secure"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func New(lc fx.Lifecycle, log *zap.Logger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	handler := gin.Default()
	handler.Use(gzip.Gzip(gzip.DefaultCompression))
	handler.Use(ginzap.Ginzap(log, time.RFC3339, true))
	handler.Use(ginzap.RecoveryWithZap(log, true))

	securityConfig := secure.DefaultConfig()
	securityConfig.SSLRedirect = false
	handler.Use(secure.New(securityConfig))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				log.Info("[Server] Failed to listen tcp at", zap.String("addr", srv.Addr))
				return err
			}
			go func() {
				err := srv.Serve(ln)
				if err != nil {
					log.Info("[Server] Failed to start HTTP Server at", zap.String("addr", srv.Addr))
				}
				log.Info("[Server] Succeeded to start HTTP Server at", zap.String("addr", srv.Addr))
			}()
			return nil

		},
		OnStop: func(ctx context.Context) error {
			err := srv.Shutdown(ctx)
			if err != nil {
				log.Info("[Server] Failed to stop HTTP Server")
			}
			log.Info("[Server] HTTP Server is stopped")
			return nil
		},
	})

	return handler
}
