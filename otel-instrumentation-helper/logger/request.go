package logger

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
	"net"
	"net/http"
	"time"
)

func SetupRequestLog(r *chi.Mux) {
	l := Logger()
	r.Use(ReqLogger("router", l))
}

// ReqLogger returns a request logging middleware
func ReqLogger(category string, logger *slog.Logger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			reqID := middleware.GetReqID(r.Context())
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now()
			defer func() {
				remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
				if err != nil {
					remoteIP = r.RemoteAddr
				}
				scheme := "http"
				if r.TLS != nil {
					scheme = "https"
				}
				span := trace.SpanFromContext(r.Context())
				fields := []slog.Attr{
					slog.Int("status_code", ww.Status()),
					slog.Int("bytes", ww.BytesWritten()),
					slog.Int64("duration", int64(time.Since(t1))),
					slog.String("duration_display", time.Since(t1).String()),
					slog.String("category", category),
					slog.String("remote_ip", remoteIP),
					slog.String("proto", r.Proto),
					slog.String("method", r.Method),
					slog.String("trace_id", span.SpanContext().TraceID().String()),
					slog.String("span_id", span.SpanContext().SpanID().String()),
					slog.String("path", r.RequestURI),
					slog.String("url", fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)),
				}
				if len(reqID) > 0 {
					fields = append(fields, slog.String("request_id", reqID))
				}
				logger.With(fields).Info("RequestReceived")
			}()

			h.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
