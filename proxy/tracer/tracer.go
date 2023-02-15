package tracer

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func InitTracer(jaegerURL string, serviceName string) (*tracesdk.TracerProvider, error) {
	exporter, err := newJaegerExporter(jaegerURL)
	if err != nil {
		return nil, fmt.Errorf("initialize exporter: %w", err)
	}

	tp, err := newTraceProvider(exporter, serviceName)
	if err != nil {
		return nil, fmt.Errorf("initialize provider: %w", err)
	}

	otel.SetTracerProvider(tp) // !!!!!!!!!!!

	return tp, nil
}

// newJaegerExporter creates new jaeger exporter
//
//	url example - http://localhost:14268/api/traces
func newJaegerExporter(URL string) (tracesdk.SpanExporter, error) {
	return jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(URL)))
}

func newTraceProvider(exp tracesdk.SpanExporter, serviceName string) (*tracesdk.TracerProvider, error) {
	// Ensure default SDK resources and the required service name are set.

	return tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			attribute.String("env", "stage"),
		)),
	), nil
}

// Cleanly shutdown and flush telemetry when the application exits.
// Do not make the application hang when it is shutdown.
func Stop(ctx context.Context, tp *tracesdk.TracerProvider) {
	ctx_timeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	if err := tp.Shutdown(ctx_timeout); err != nil {
		log.Fatal(err)
	}
}
