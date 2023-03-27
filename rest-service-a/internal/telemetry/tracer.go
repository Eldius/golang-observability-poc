package telemetry

import (
	"context"
	"github.com/eldius/rest-api/internal/config"
	"github.com/rs/zerolog/log"
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

func initTracer(ctx context.Context) {
	log.Debug().Msg("init tracer begin")

	// initialize trace provider
	tp := initTracerProvider(ctx)
	// set global tracer provider & text propagators

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	log.Debug().Msg("finished tracer configuration")

	log.Debug().Msg("ending tracer provider")

	go wait(tp)
}

func initTracerProvider(ctx context.Context) *sdktrace.TracerProvider {
	exporter := chooseTraceExporter(ctx)

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(config.GetServiceName()),
		semconv.ServiceVersionKey.String("v0.0.0"),
		attribute.String("environment", config.GetEnvironment()),
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

func stdoutTraceExporter() sdktrace.SpanExporter {
	log.Debug().Msg("configuring stdout trace export")

	var err error
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint(), stdouttrace.WithWriter(log.Logger))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to setup exporter")
	}

	return exporter
}

func otelTraceExporter(ctx context.Context) sdktrace.SpanExporter {
	log.Debug().Msgf("configuring trace export for '%s'", config.GetOtelTraceEndpoint())

	var err error
	conn, err := grpc.DialContext(
		ctx,
		config.GetOtelTraceEndpoint(),
		// Note the use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create gRPC connection to collector")
	}

	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to setup exporter")
	}

	return exporter
}

func chooseTraceExporter(ctx context.Context) sdktrace.SpanExporter {
	endpoint := config.GetOtelTraceEndpoint()
	log.Debug().Msgf("otel_endpoint: %s", endpoint)
	if endpoint == "" {
		return stdoutTraceExporter()
	} else {
		return otelTraceExporter(ctx)
	}
}

func wait(tp *sdktrace.TracerProvider) {
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Debug().Err(err).Msg("error shutting down tracer provider")
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
