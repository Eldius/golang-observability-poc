package logger

import (
	"context"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

func GetLogger(ctx context.Context) zerolog.Logger {
	span := trace.SpanFromContext(ctx)

	return zerolog.Ctx(ctx).
		With().
		Str("trace_id", span.SpanContext().TraceID().String()).
		Str("span_id", span.SpanContext().SpanID().String()).
		Logger().Level(zerolog.GlobalLevel())
}
