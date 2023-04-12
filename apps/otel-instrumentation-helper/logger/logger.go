package logger

import (
	"context"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"go.opentelemetry.io/otel/trace"
)

var serviceName string

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

func SetupLogs(level, service string) {
	logLevel := strings.ToLower(level)
	switch logLevel {
	case zerolog.LevelPanicValue:
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case zerolog.LevelFatalValue:
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case zerolog.LevelErrorValue:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case zerolog.LevelWarnValue:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case zerolog.LevelInfoValue:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case zerolog.LevelDebugValue:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case zerolog.LevelTraceValue:
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	serviceName = service

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack //nolint:reassign // setting stack traces marshaller

	log.Info().Str("setup_log_level", zerolog.GlobalLevel().String()).Msg("SetupLogsEnd")
}
