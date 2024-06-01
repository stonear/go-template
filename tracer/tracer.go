package tracer

import (
	"context"
	"time"

	"github.com/stonear/go-template/config"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Load(lc fx.Lifecycle, config *config.Config, log *otelzap.Logger) {
	if !config.EnableTelemetry {
		log.Info("[Tracer] Telemetry is disabled")
		return
	}

	// Configure OpenTelemetry with sensible defaults.
	uptrace.ConfigureOpentelemetry(
		uptrace.WithDSN(config.UptraceDsn),
		uptrace.WithServiceName(config.AppName),
		uptrace.WithServiceVersion(config.AppVersion),
	)

	err := runtime.Start(runtime.WithMinimumReadMemStatsInterval(15 * time.Second))
	if err != nil {
		log.Fatal("[Tracer] Unable to start runtime instrumentation", zap.Error(err))
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			if !config.EnableTelemetry {
				return nil
			}

			// Send buffered spans and free resources.
			return uptrace.Shutdown(ctx)
		},
	})
}
