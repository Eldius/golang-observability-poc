package logger

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/trace"
	"io"
	"log/slog"
	"os"
	"slices"
	"strings"
)

var (
	logKeys = []string{
		"host",
		"level",
		"message",
		"time",
		"error",
		"source",
		"function",
		"file",
		"line",
		"trace_id",
		"request",
	}

	logger *slog.Logger
)

func GetLogger(ctx context.Context) *slog.Logger {
	span := trace.SpanFromContext(ctx)

	return logger.With(
		slog.String("trace_id", span.SpanContext().TraceID().String()),
		slog.String("span_id", span.SpanContext().SpanID().String()),
	)
}

func Logger() *slog.Logger {
	return logger
}

func SetupLogs(logLevel, service string) error {
	var h slog.Handler
	var w io.Writer = os.Stdout

	replaceAttrFunc := func(groups []string, a slog.Attr) slog.Attr {
		if slices.Contains(logKeys, a.Key) {
			return a
		}
		if strings.HasPrefix(a.Key, "request.") || strings.HasPrefix(a.Key, "response.") || strings.HasPrefix(a.Key, "service.") {
			return a
		}
		if a.Key == slog.MessageKey {
			a.Key = "message"
			return a
		}
		a.Key = fmt.Sprintf("custom.%s.%s", service, a.Key)
		return a
	}

	h = slog.NewJSONHandler(w, &slog.HandlerOptions{
		AddSource:   true,
		Level:       parseLogLevel(logLevel),
		ReplaceAttr: replaceAttrFunc,
	})
	logger = slog.New(h)
	host, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	slog.SetDefault(logger.With(
		slog.String("service.name", service),
		slog.String("host", host),
	))

	return nil
}

func parseLogLevel(logLevel string) slog.Level {
	switch strings.ToLower(logLevel) {
	case "error":
		return slog.LevelError
	case "warn":
		return slog.LevelWarn
	case "info":

		return slog.LevelInfo
	case "debug":
		return slog.LevelDebug
	default:
		return slog.LevelInfo
	}
}
