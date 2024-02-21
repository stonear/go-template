package tracer

import (
	"context"
	"os"

	"github.com/uptrace/uptrace-go/uptrace"
	"go.uber.org/fx"
)

func Load(lc fx.Lifecycle) {
	// Configure OpenTelemetry with sensible defaults.
	uptrace.ConfigureOpentelemetry(
		// copy your project DSN here or use UPTRACE_DSN env var
		// uptrace.WithDSN("https://token@api.uptrace.dev/project_id"),

		uptrace.WithServiceName(os.Getenv("APP_NAME")),
		uptrace.WithServiceVersion(os.Getenv("APP_VERSION")),
	)

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			// Send buffered spans and free resources.
			return uptrace.Shutdown(ctx)
		},
	})
}
