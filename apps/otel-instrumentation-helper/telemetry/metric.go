package telemetry

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/eldius/golang-observability-poc/apps/otel-instrumentation-helper/logger"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
	"google.golang.org/grpc/encoding/gzip"
)

func initMetrics(opt *options) {
	l := logger.Logger()
	l.Debug().Msg("init tracer begin")

	// initialize trace provider
	mp := initMetricsProvider(opt)
	// set global tracer provider & text propagators

	global.SetMeterProvider(mp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	l.Debug().Msg("finished metrics provider configuration")

	l.Debug().Msg("starting runtime instrumentation")
	if err := runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second)); err != nil {
		l.Warn().Err(err).Msg("failed to start runtime monitoring")
		return
	}

	l.Debug().Msg("ending metrics provider")

	go waitMetrics(mp)
}

func initMetricsProvider(opt *options) otelmetric.MeterProvider {
	if opt.metricsEndpoint == "" {
		return otelmetric.NewNoopMeterProvider()
	}
	exporter := otelMetricsExporter(opt)

	provider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(exporter)),
		metric.WithResource(defaultResources(opt)))

	global.SetMeterProvider(provider)

	return provider
}

func defaultResources(opt *options) *resource.Resource {
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(opt.serviceName),
		semconv.ServiceVersionKey.String(opt.version),
		attribute.String("environment", opt.env),
	)
	return res
}

func otelMetricsExporter(opt *options) metric.Exporter {
	l := logger.Logger()
	l.Debug().Msgf("configuring metric export for '%s'", opt.metricsEndpoint)

	var opts []otlpmetricgrpc.Option

	opts = append(opts,
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(opt.metricsEndpoint),
		otlpmetricgrpc.WithCompressor(gzip.Name),
		otlpmetricgrpc.WithTimeout(10*time.Second))

	exporter, err := otlpmetricgrpc.New(
		opt.ctx,
		opts...,
	)
	if err != nil {
		l.Warn().Err(err).Msg("failed to configure otel metrics exporter")
		return nil
	}

	return exporter
}

func waitMetrics(mp otelmetric.MeterProvider) {
	defer func() {
		if p, ok := mp.(*metric.MeterProvider); ok {
			l := logger.Logger()
			if err := p.Shutdown(context.Background()); err != nil {
				l.Debug().Err(err).Msg("error shutting down metric provider")
			}
		}
	}()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	<-ctx.Done()
}
