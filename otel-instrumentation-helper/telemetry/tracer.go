package telemetry

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/eldius/golang-observability-poc/otel-instrumentation-helper/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var Tracer trace.Tracer

func initTracer(opt *options) {
	l := logger.Logger()
	l.Debug("init tracer begin")

	// initialize trace provider
	tp := initTracerProvider(opt)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	l.Debug("finished tracer configuration")

	l.Debug("ending tracer provider")

	go waitTraces(tp)
}

func initTracerProvider(opt *options) trace.TracerProvider {
	if opt.tracesEndpoint == "" {
		return nil
	}
	exporter := otelTraceExporter(opt)

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(opt.serviceName),
		semconv.ServiceVersionKey.String(opt.version),
		attribute.String("environment", opt.env),
	)

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	bsp := sdktrace.NewBatchSpanProcessor(exporter)
	provider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	// set global tracer provider & text propagators
	otel.SetTracerProvider(provider)
	Tracer = provider.Tracer(opt.serviceName)

	return provider
}

func otelTraceExporter(opt *options) sdktrace.SpanExporter {
	l := logger.Logger()
	l.With(slog.String("metrics.endpoint", opt.metricsEndpoint)).Debug("configuring trace exporter")

	var err error
	conn, err := grpc.DialContext(
		opt.ctx,
		opt.tracesEndpoint,
		// Note the use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		l.With("error", err).Error("failed to create gRPC connection to collector")
		panic(err)
	}

	exporter, err := otlptracegrpc.New(opt.ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		l.With("error", err).Error("failed to setup exporter")
		panic(err)
	}

	return exporter
}

func waitTraces(tp trace.TracerProvider) {
	defer func() {
		if p, ok := tp.(*sdktrace.TracerProvider); ok {
			l := logger.Logger()
			if err := p.Shutdown(context.Background()); err != nil {
				l.With("error", err).Debug("error shutting down tracer provider")
			}
		}
	}()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	<-ctx.Done()
}

func NotifyError(ctx context.Context, err error) {
	span := trace.SpanFromContext(ctx)
	span.RecordError(err, trace.WithStackTrace(true))
}

func AddTraceAttributes(ctx context.Context, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attrs...)
}
