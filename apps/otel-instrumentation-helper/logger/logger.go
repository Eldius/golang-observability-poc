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
	case zerolog.LevelPanicValue:
		fmt.Println("panic")
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case zerolog.LevelFatalValue:
		fmt.Println("fatal")
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case zerolog.LevelErrorValue:
		fmt.Println("error")
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case zerolog.LevelWarnValue:
		fmt.Println("warn")
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case zerolog.LevelInfoValue:
		fmt.Println("info")
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case zerolog.LevelDebugValue:
		fmt.Println("debug")
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case zerolog.LevelTraceValue:
		fmt.Println("trace")
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	default:
		fmt.Println("default")
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	serviceName = service
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	log.Info().Str("setup_log_level", zerolog.GlobalLevel().String()).Msg("SetupLogsEnd")
	fmt.Println("global log level:", zerolog.GlobalLevel().String())
}
