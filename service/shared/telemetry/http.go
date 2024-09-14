package telemetry

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/semconv/v1.13.0/httpconv"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

const tracerKey = "edot-custom-tracer"

func HttpOtel(service string) echo.MiddlewareFunc {

	tracerProvider := otel.GetTracerProvider()

	tracer := tracerProvider.Tracer(
		"edot-custom-tracer",
		oteltrace.WithInstrumentationVersion("1.0.0"),
	)

	propagators := otel.GetTextMapPropagator()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Path() == "/metrics" {
				return next(c)
			}

			c.Set(tracerKey, tracer)
			request := c.Request()
			savedCtx := request.Context()

			defer func() {
				request = request.WithContext(savedCtx)
				c.SetRequest(request)
			}()

			ctx := propagators.Extract(savedCtx, propagation.HeaderCarrier(request.Header))

			opts := []oteltrace.SpanStartOption{
				oteltrace.WithAttributes(httpconv.ServerRequest(service, request)...),
				oteltrace.WithSpanKind(oteltrace.SpanKindServer),
				oteltrace.WithAttributes(attribute.String(echo.HeaderXRequestID, request.Header.Get(echo.HeaderXRequestID))),
			}

			if path := c.Path(); path != "" {
				rAttr := semconv.HTTPRoute(path)
				opts = append(opts, oteltrace.WithAttributes(rAttr))
			}

			spanName := c.Path()
			if spanName == "" {
				spanName = fmt.Sprintf("HTTP %s route not found", request.Method)
			}

			ctx, span := tracer.Start(ctx, spanName, opts...)
			defer span.End()

			// pass the span through the request context
			c.SetRequest(request.WithContext(ctx))

			// serve the request to the next middleware
			err := next(c)
			if err != nil {
				span.SetAttributes(attribute.String("echo.error", err.Error()))
				// invokes the registered HTTP error handler
				c.Error(err)
			}

			status := c.Response().Status
			span.SetStatus(httpconv.ServerStatus(status))
			if status > 0 {
				span.SetAttributes(semconv.HTTPStatusCode(status))
			}

			return err
		}
	}
}
