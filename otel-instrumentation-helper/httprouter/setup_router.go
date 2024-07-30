package httprouter

import (
	"github.com/eldius/golang-observability-poc/otel-instrumentation-helper/logger"
	"github.com/eldius/golang-observability-poc/otel-instrumentation-helper/telemetry"
	"net/http"
)

func SetupRouter(serviceName string, r http.Handler) http.Handler {
	return telemetry.TraceMiddleware(serviceName)(logger.SetupRequestLog(r))
}
