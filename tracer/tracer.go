package tracer

import (
	"context"
	"os"
	"time"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Load(lc fx.Lifecycle, log *otelzap.Logger) {
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
			// Send buffered spans and free resources.
			return uptrace.Shutdown(ctx)
		},
	})
}
