package infra

import (
	"log"
	"shared/telemetry/trace"
	"shared/telemetry/trace/exporter"
	"sync"
)

var (
	tracerProviderOnce sync.Once
	tracerProvider     *trace.TraceProvider
)

func LoadTraceProvider() *trace.TraceProvider {
	tracerProviderOnce.Do(func() {

		var err error

		tracerProvider, err = trace.New(
			"product",
			exporter.NewOTLP(LoadConfig().Telemetry.Otlp),
		)
		if err != nil {
			log.Fatal(err)
		}

	})

	return tracerProvider
}
