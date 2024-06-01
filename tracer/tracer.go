package tracer

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Load(lc fx.Lifecycle, log *otelzap.Logger) {
	enableTelemetryStr := os.Getenv("ENABLE_TELEMETRY")
	enableTelemetry, _ := strconv.ParseBool(enableTelemetryStr)
	if !enableTelemetry {
		return
	}

	// Configure OpenTelemetry with sensible defaults.
	uptrace.ConfigureOpentelemetry(
		// copy your project DSN here or use UPTRACE_DSN env var
		// uptrace.WithDSN("https://token@api.uptrace.dev/project_id"),

		uptrace.WithServiceName(os.Getenv("APP_NAME")),
		uptrace.WithServiceVersion(os.Getenv("APP_VERSION")),
	)

	err := runtime.Start(runtime.WithMinimumReadMemStatsInterval(15 * time.Second))
	if err != nil {
		log.Fatal("[Tracer] Unable to start runtime instrumentation", zap.Error(err))
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			if !enableTelemetry {
				return nil
			}

			// Send buffered spans and free resources.
			return uptrace.Shutdown(ctx)
		},
	})
}
