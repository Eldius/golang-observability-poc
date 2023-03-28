package telemetry

import (
	"context"
	"encoding/json"
	"github.com/eldius/rest-api/internal/config"
	"github.com/eldius/rest-api/internal/logger"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
	"google.golang.org/grpc/encoding/gzip"
	"io"
	"os"
	"os/signal"
	"time"
)

func initMetrics(ctx context.Context) {
	log.Debug().Msg("init tracer begin")

	// initialize trace provider
	mp := initMetricsProvider(ctx)
	// set global tracer provider & text propagators

	global.SetMeterProvider(mp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	log.Debug().Msg("finished metrics provider configuration")

	log.Debug().Msg("starting runtime instrumentation")
	if err := runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second)); err != nil {
		log.Warn().Err(err).Msg("failed to start runtime monitoring")
		return
	}

	log.Debug().Msg("ending metrics provider")

	go waitMetrics(mp)
}

func initMetricsProvider(ctx context.Context) *metric.MeterProvider {
	exporter := chooseMetricsExporter(ctx)

	provider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(exporter)),
		metric.WithResource(defaultResources()))

	//provider := metric.NewMeterProvider(
	//	metric.WithReader(metric.NewPeriodicReader(exporter)),
	//	metric.WithResource(defaultResources()))

	global.SetMeterProvider(provider)

	return provider
}

func defaultResources() *resource.Resource {
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(config.GetServiceName()),
		semconv.ServiceVersionKey.String(config.GetVersion()),
		attribute.String("environment", config.GetEnvironment()),
	)
	return res
}

type otelStdoutExporterWriter struct {
	l zerolog.Logger
	t string
}

func (e *otelStdoutExporterWriter) Write(p []byte) (n int, err error) {
	e.l.Debug().Str("type", e.t).Msg(string(p))
	return len(p), nil
}

func newStdoutWriter(t string) io.Writer {
	return &otelStdoutExporterWriter{
		l: logger.GetLogger(context.Background()),
		t: t,
	}
}

func stdoutMetricsExporter() metric.Exporter {
	enc := json.NewEncoder(newStdoutWriter("otel_metrics"))
	//enc.SetIndent("", "  ")
	exp, err := stdoutmetric.New(stdoutmetric.WithEncoder(enc))
	if err != nil {
		panic(err)
	}
	return exp
}

func chooseMetricsExporter(ctx context.Context) metric.Exporter {
	endpoint := config.GetOtelMetricsEndpoint()
	log.Debug().Msgf("otel_metrics_endpoint: %s", endpoint)
	if endpoint == "" {
		return stdoutMetricsExporter()
	} else {
		return otelMetricsExporter(ctx)
	}
}

func otelMetricsExporter(ctx context.Context) metric.Exporter {
	log.Debug().Msgf("configuring metric export for '%s'", config.GetOtelMetricsEndpoint())

	//var err error
	//conn, err := grpc.DialContext(
	//	ctx,
	//	config.GetOtelTraceEndpoint(),
	//	grpc.WithTransportCredentials(insecure.NewCredentials()),
	//	grpc.WithBlock(),
	//)
	//if err != nil {
	//	log.Fatal().Err(err).Msg("failed to create gRPC connection to collector")
	//}

	var opts []otlpmetricgrpc.Option
	opts = append(opts, otlpmetricgrpc.WithInsecure())
	//opts = append(opts, otlpmetricgrpc.WithGRPCConn(conn))
	opts = append(opts, otlpmetricgrpc.WithEndpoint(config.GetOtelMetricsEndpoint()))
	opts = append(opts, otlpmetricgrpc.WithCompressor(gzip.Name))
	opts = append(opts, otlpmetricgrpc.WithTimeout(10*time.Second))

	exporter, err := otlpmetricgrpc.New(
		ctx,
		opts...,
	)
	if err != nil {
		log.Warn().Err(err).Msg("failed to configure otel metrics exporter")
		return nil
	}

	return exporter
}

func waitMetrics(mp *metric.MeterProvider) {
	defer func() {
		if err := mp.Shutdown(context.Background()); err != nil {
			log.Debug().Err(err).Msg("error shutting down metric provider")
		}
	}()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	<-ctx.Done()

}
