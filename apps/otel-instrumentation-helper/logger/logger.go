package logger

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"go.opentelemetry.io/otel/trace"
	"strings"
)

var (
	serviceName string
)

func GetLogger(ctx context.Context) zerolog.Logger {
	span := trace.SpanFromContext(ctx)

	return zerolog.Ctx(ctx).
		With().
		Caller().
		Str("trace_id", span.SpanContext().TraceID().String()).
		Str("span_id", span.SpanContext().SpanID().String()).
		Logger().Level(zerolog.GlobalLevel())
}

func Logger() zerolog.Logger {
	return zerolog.Ctx(context.Background()).
		With().
		Caller().
		Str("service", serviceName).
		Logger().Level(zerolog.GlobalLevel())
}

func SetupLogs(level string, service string) {
	logLevel := strings.ToLower(level)
	fmt.Printf("configuring log level: '%s'\n", logLevel)
	switch logLevel {
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	serviceName = service
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	log.Info().Str("setup_log_level", zerolog.GlobalLevel().String())
}
