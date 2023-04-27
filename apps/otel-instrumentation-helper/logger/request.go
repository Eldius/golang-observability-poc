package logger

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
	"net"
	"net/http"
	"time"
)

func SetupRequestLog(r *chi.Mux) {
	l := Logger()
	r.Use(ReqLogger("router", l))
}

// ReqLogger returns a request logging middleware
func ReqLogger(category string, logger logrus.FieldLogger) func(h http.Handler) http.Handler {
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
				fields := logrus.Fields{
					"status_code":      ww.Status(),
					"bytes":            ww.BytesWritten(),
					"duration":         int64(time.Since(t1)),
					"duration_display": time.Since(t1).String(),
					"category":         category,
					"remote_ip":        remoteIP,
					"proto":            r.Proto,
					"method":           r.Method,
					"trace_id":         span.SpanContext().TraceID().String(),
					"span_id":          span.SpanContext().SpanID().String(),
				}
				if len(reqID) > 0 {
					fields["request_id"] = reqID
				}
				logger.WithFields(fields).Infof("%s://%s%s", scheme, r.Host, r.RequestURI)
			}()

			h.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
