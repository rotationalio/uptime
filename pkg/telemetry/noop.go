package telemetry

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/log/global"
	nooplog "go.opentelemetry.io/otel/log/noop"
	noopmeter "go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	nooptrace "go.opentelemetry.io/otel/trace/noop"
)

func disableTelemetry(ctx context.Context) {
	var err error
	if upResource, err = resource.New(ctx); err != nil {
		initerr = errors.Join(initerr, err)
	}

	upPropagator = propagation.TraceContext{}
	otel.SetTextMapPropagator(upPropagator)

	otel.SetTracerProvider(nooptrace.NewTracerProvider())
	otel.SetMeterProvider(noopmeter.NewMeterProvider())
	global.SetLoggerProvider(nooplog.NewLoggerProvider())

	disabled = true

	// Setup application logging with console only logging.
	if err = initializeLogging(); err != nil {
		initerr = errors.Join(initerr, err)
	}
}
