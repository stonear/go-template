package tracer

import (
	"context"
	"os"
	"time"

	"github.com/stonear/go-template/config"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Load(lc fx.Lifecycle, config *config.Config, log *otelzap.Logger) {
	if !config.EnableTelemetry {
		log.Info("[Tracer] Telemetry is disabled")
		return
	}

	// Set environment variables for OTLP exporter
	os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", config.OtelExporterOtlpEndpoint)
	os.Setenv("OTEL_EXPORTER_OTLP_HEADERS", config.OtelExporterOtlpHeaders)
	os.Setenv("OTEL_SERVICE_NAME", config.AppName)

	// Create a new OTLP trace exporter
	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(),
	)
	if err != nil {
		log.Fatal("[Tracer] Failed to create OTLP trace exporter", zap.Error(err))
	}

	// Create a new resource with the service name and version
	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceName(config.AppName),
			semconv.ServiceVersion(config.AppVersion),
		),
	)
	if err != nil {
		log.Fatal("[Tracer] Failed to create resource", zap.Error(err))
	}

	// Create a new trace provider with the exporter and resource
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	// Set the global trace provider
	otel.SetTracerProvider(tp)

	// Set the global propagator to trace context (for context propagation)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Start runtime instrumentation
	err = runtime.Start(runtime.WithMinimumReadMemStatsInterval(15 * time.Second))
	if err != nil {
		log.Fatal("[Tracer] Unable to start runtime instrumentation", zap.Error(err))
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			if !config.EnableTelemetry {
				return nil
			}

			// Shutdown the trace provider to flush any spans and free resources
			return tp.Shutdown(ctx)
		},
	})
}
