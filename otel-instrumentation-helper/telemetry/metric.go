package telemetry

import (
    "context"
    "encoding/json"
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

func initMetrics(opt *options) {
    log.Debug().Msg("init tracer begin")

    // initialize trace provider
    mp := initMetricsProvider(opt)
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

func initMetricsProvider(opt *options) *metric.MeterProvider {
    exporter := chooseMetricsExporter(opt)

    provider := metric.NewMeterProvider(
        metric.WithReader(metric.NewPeriodicReader(exporter)),
        metric.WithResource(defaultResources(opt)))

    //provider := metric.NewMeterProvider(
    //	metric.WithReader(metric.NewPeriodicReader(exporter)),
    //	metric.WithResource(defaultResources()))

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

type otelStdoutExporterWriter struct {
    l zerolog.Logger
    t string
}

func (e *otelStdoutExporterWriter) Write(p []byte) (n int, err error) {
    e.l.Debug().Str("type", e.t).Msg(string(p))
    return len(p), nil
}

func newStdoutWriter(t string, opt *options) io.Writer {
    return &otelStdoutExporterWriter{
        l: zerolog.Ctx(opt.ctx).
            With().
            Logger().Level(zerolog.GlobalLevel()),
        t: t,
    }
}

func stdoutMetricsExporter(opt *options) metric.Exporter {
    enc := json.NewEncoder(newStdoutWriter("otel_metrics", opt))
    //enc.SetIndent("", "  ")
    exp, err := stdoutmetric.New(stdoutmetric.WithEncoder(enc))
    if err != nil {
        panic(err)
    }
    return exp
}

func chooseMetricsExporter(opt *options) metric.Exporter {
    endpoint := opt.metricsEndpoint
    log.Debug().Msgf("otel_metrics_endpoint: %s", endpoint)
    if endpoint == "" {
        return stdoutMetricsExporter(opt)
    } else {
        return otelMetricsExporter(opt)
    }
}

func otelMetricsExporter(opt *options) metric.Exporter {
    log.Debug().Msgf("configuring metric export for '%s'", opt.metricsEndpoint)

    var opts []otlpmetricgrpc.Option
    opts = append(opts, otlpmetricgrpc.WithInsecure())
    //opts = append(opts, otlpmetricgrpc.WithGRPCConn(conn))
    opts = append(opts, otlpmetricgrpc.WithEndpoint(opt.metricsEndpoint))
    opts = append(opts, otlpmetricgrpc.WithCompressor(gzip.Name))
    opts = append(opts, otlpmetricgrpc.WithTimeout(10*time.Second))

    exporter, err := otlpmetricgrpc.New(
        opt.ctx,
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