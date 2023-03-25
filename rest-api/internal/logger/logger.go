package logger

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/trace"
)

type Logger struct {
	l   *zerolog.Logger
	ctx context.Context
}

func GetLogger(ctx context.Context) *Logger {
	return &Logger{
		l:   log.Ctx(ctx),
		ctx: ctx,
	}
}

func (l *Logger) Panic(format string, v ...interface{}) {
	span := trace.SpanFromContext(l.ctx)
	l.l.Panic().
		Str("trace_id", span.SpanContext().TraceID().String()).
		Str("span_id", span.SpanContext().SpanID().String()).
		Msgf(format, v)
}

func (l *Logger) Fatal(format string, v ...interface{}) {
	span := trace.SpanFromContext(l.ctx)
	l.l.Fatal().
		Str("trace_id", span.SpanContext().TraceID().String()).
		Str("span_id", span.SpanContext().SpanID().String()).
		Msgf(format, v)
}

func (l *Logger) Error(format string, v ...interface{}) {
	span := trace.SpanFromContext(l.ctx)
	l.l.Error().
		Str("trace_id", span.SpanContext().TraceID().String()).
		Str("span_id", span.SpanContext().SpanID().String()).
		Msgf(format, v)
}

func (l *Logger) Warn(format string, v ...interface{}) {
	span := trace.SpanFromContext(l.ctx)
	l.l.Warn().
		Str("trace_id", span.SpanContext().TraceID().String()).
		Str("span_id", span.SpanContext().SpanID().String()).
		Msgf(format, v)
}

func (l *Logger) Info(format string, v ...interface{}) {
	span := trace.SpanFromContext(l.ctx)
	l.l.Info().
		Str("trace_id", span.SpanContext().TraceID().String()).
		Str("span_id", span.SpanContext().SpanID().String()).
		Msgf(format, v)
}

func (l *Logger) Debug(format string, v ...interface{}) {
	span := trace.SpanFromContext(l.ctx)
	l.l.Debug().
		Str("trace_id", span.SpanContext().TraceID().String()).
		Str("span_id", span.SpanContext().SpanID().String()).
		Msgf(format, v)
}

func (l *Logger) Trace(format string, v ...interface{}) {
	span := trace.SpanFromContext(l.ctx)
	l.l.Trace().
		Str("trace_id", span.SpanContext().TraceID().String()).
		Str("span_id", span.SpanContext().SpanID().String()).
		Msgf(format, v)
}
