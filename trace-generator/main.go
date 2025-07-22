package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func main() {
	ctx := context.Background()

	// Create OTLP GRPC exporter
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint("localhost:4317"),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Create a single resource - we'll set service names per span
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("trace-generator"),
			semconv.ServiceVersionKey.String("1.0.0"),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Create single trace provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	defer tp.Shutdown(ctx)

	tracer := tp.Tracer("trace-generator")

	// Generate test traces with proper parent-child relationships
	generateTestTraces(ctx, tracer)
}

func generateTestTraces(ctx context.Context, tracer trace.Tracer) {
	// Generate traces that match the expected structure:
	// root_span (appserver) -> getData (appserver) -> Signal/read spans (datasource-proxy)
	
	for traceNum := 1; traceNum <= 3; traceNum++ {
		// Create a new trace for each iteration
		rootCtx, rootSpan := tracer.Start(ctx, "handle_request")
		rootSpan.SetAttributes(
			attribute.String("service.name", "appserver"),
		)
		
		// Create the single getData operation per trace (this is critical!)
		getDataCtx, getDataSpan := tracer.Start(rootCtx, "compute.v1.ComputeEngine/getData")
		getDataSpan.SetAttributes(
			attribute.String("service.name", "appserver"),
		)
		
		// Create multiple Signal/read spans as descendants of getData
		numReads := 2 + traceNum // Variable number of reads per trace
		for i := 0; i < numReads; i++ {
			_, readSpan := tracer.Start(getDataCtx, "data.signal.v1.Signal/read")
			readSpan.SetAttributes(
				attribute.String("service.name", "datasource-proxy"),
				attribute.String("signal.id", fmt.Sprintf("signal_%d_%d", traceNum, i)),
			)
			
			// Simulate different read durations that should be summed
			readDuration := time.Duration(50+i*25) * time.Millisecond
			time.Sleep(readDuration)
			readSpan.End()
		}
		
		// End getData span
		time.Sleep(10 * time.Millisecond) // Small additional work
		getDataSpan.End()
		
		// End root span
		time.Sleep(5 * time.Millisecond)
		rootSpan.End()
		
		fmt.Printf("Generated trace %d with %d Signal/read spans\n", traceNum, numReads)
	}
	
	fmt.Printf("Generated test traces matching field summing collector assumptions\n")
}