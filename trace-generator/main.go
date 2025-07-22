package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func main() {
	ctx := context.Background()

	// Create OTLP HTTP exporter
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint("http://localhost:4318"),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Create resource
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("trace-generator"),
			semconv.ServiceVersionKey.String("1.0.0"),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Create trace provider
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)
	defer tp.Shutdown(ctx)

	otel.SetTracerProvider(tp)
	tracer := tp.Tracer("trace-generator")

	// Generate test traces with various field patterns for summing
	generateTestTraces(ctx, tracer)
}

func generateTestTraces(ctx context.Context, tracer trace.Tracer) {
	// Generate traces with numeric fields that should be summed
	testCases := []struct {
		name     string
		duration time.Duration
		count    int
		value    float64
	}{
		{"fast-operation", 100 * time.Millisecond, 5, 10.5},
		{"medium-operation", 500 * time.Millisecond, 3, 25.0},
		{"slow-operation", 1 * time.Second, 2, 50.5},
	}

	for _, tc := range testCases {
		for i := 0; i < tc.count; i++ {
			_, span := tracer.Start(ctx, tc.name)
			
			span.SetAttributes(
				attribute.String("operation.type", tc.name),
				attribute.Int("operation.count", 1),
				attribute.Float64("operation.value", tc.value),
				attribute.Int("iteration", i+1),
			)

			// Simulate work
			time.Sleep(tc.duration)
			span.End()
		}
	}

	fmt.Printf("Generated test traces with summing fields\n")
}