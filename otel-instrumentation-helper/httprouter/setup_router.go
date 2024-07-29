package httprouter

import (
	"github.com/eldius/golang-observability-poc/otel-instrumentation-helper/logger"
	"github.com/eldius/golang-observability-poc/otel-instrumentation-helper/telemetry"
	"net/http"
)

func SetupRouter(r http.Handler) http.Handler {
	router := logger.SetupRequestLog(r)
	return telemetry.TracedRouter(router)
}
