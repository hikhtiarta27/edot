package exporter

import (
	"context"
	"log"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
)

func NewOTLP(endpoint string) *otlptrace.Exporter {
	ctx := context.Background()
	traceClient := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(endpoint),
	)
	traceExp, err := otlptrace.New(ctx, traceClient)
	if err != nil {
		log.Fatal(err, "Failed to create the collector trace exporter")
	}

	return traceExp
}
