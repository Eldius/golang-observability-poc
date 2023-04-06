package telemetry

import (
	"context"
	"github.com/eldius/golang-observability-poc/apps/otel-instrumentation-helper/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"os/signal"
)

var Tracer trace.Tracer

func initTracer(opt *options) {
	l := logger.Logger()
	l.Debug().Msg("init tracer begin")

	// initialize trace provider
	tp := initTracerProvider(opt)
	// set global tracer provider & text propagators

	Tracer = tp.Tracer(opt.serviceName)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	l.Debug().Msg("finished tracer configuration")

	l.Debug().Msg("ending tracer provider")

	go waitTraces(tp)
}

func initTracerProvider(opt *options) trace.TracerProvider {
	if opt.tracesEndpoint == "" {
		return trace.NewNoopTracerProvider()
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

	otel.SetTracerProvider(provider)

	return provider
}

func stdoutTraceExporter(opt *options) sdktrace.SpanExporter {
	l := logger.Logger()
	l.Debug().Msg("configuring stdout trace export")

	var err error
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint(), stdouttrace.WithWriter(newStdoutWriter("otel_tracing", opt)))
	if err != nil {
		l.Fatal().Err(err).Msg("failed to setup exporter")
	}

	return exporter
}

func otelTraceExporter(opt *options) sdktrace.SpanExporter {
	l := logger.Logger()
	l.Debug().Msgf("configuring trace export for '%s'", opt.tracesEndpoint)

	var err error
	conn, err := grpc.DialContext(
		opt.ctx,
		opt.tracesEndpoint,
		// Note the use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		l.Fatal().Err(err).Msg("failed to create gRPC connection to collector")
	}

	exporter, err := otlptracegrpc.New(opt.ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		l.Fatal().Err(err).Msg("failed to setup exporter")
	}

	return exporter
}

func waitTraces(tp trace.TracerProvider) {
	defer func() {
		if p, ok := tp.(*sdktrace.TracerProvider); ok {
			l := logger.Logger()
			if err := p.Shutdown(context.Background()); err != nil {
				l.Debug().Err(err).Msg("error shutting down tracer provider")
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
