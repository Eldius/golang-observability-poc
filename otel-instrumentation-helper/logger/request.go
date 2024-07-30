package logger

import (
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/otel/trace"
	"net"
	"net/http"
	"time"
)

func SetupRequestLog(r http.Handler) http.Handler {
	return ReqLogger("router")(r)
}

// ReqLogger returns a request logging middleware
func ReqLogger(category string) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			reqID := middleware.GetReqID(r.Context())
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now()
			logger := GetLogger(r.Context())
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

				reqFields := map[string]interface{}{
					"status_code":              ww.Status(),
					"bytes":                    ww.BytesWritten(),
					"duration":                 int64(time.Since(t1)),
					"request.duration_display": time.Since(t1).String(),
					"category":                 category,
					"remote_ip":                remoteIP,
					"proto":                    r.Proto,
					"method":                   r.Method,
					"trace_id":                 span.SpanContext().TraceID().String(),
					"span_id":                  span.SpanContext().SpanID().String(),
					"path":                     r.RequestURI,
					"url":                      fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI),
					"headers":                  r.Header,
					"response": map[string]interface{}{
						"headers": ww.Header(),
					},
				}
				if len(reqID) > 0 {
					reqFields["id"] = reqID
				}
				logger.With("request", reqFields).Info("RequestReceived")
			}()

			h.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
