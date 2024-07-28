package api

import (
	"encoding/json"
	"fmt"
	"github.com/eldius/golang-observability-poc/otel-instrumentation-helper/logger"
	"github.com/eldius/golang-observability-poc/otel-instrumentation-helper/telemetry"
	"github.com/eldius/golang-observability-poc/rest-service-b/internal/weather"
	"log/slog"
	"net/http"
	"time"
)

func Start(port int) {

	l := logger.Logger()

	r := http.NewServeMux()

	telemetry.AddRouteHandler(r, "/", homeHandlerfunc)
	telemetry.AddRouteHandler(r, "/health", healthHandlerfunc)
	telemetry.AddRouteHandler(r, "/weather", weatherHandlerFunc)

	logger.SetupRequestLog(r)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           r,
		ReadHeaderTimeout: 100 * time.Millisecond,
	}
	l.With(slog.String("addr", srv.Addr)).
		Info("starting")
	if err := srv.ListenAndServe(); err != nil {
		l.With("error", err).Error("filed to start server")
		panic(err)
	}
}

func homeHandlerfunc(w http.ResponseWriter, r *http.Request) {

	l := logger.GetLogger(r.Context())
	l.Info("get root begin")
	l.Info("testing: %s", r.Context())

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("home")) //nolint:errcheck // ignoring error

	l.Info("get root end")
}

func healthHandlerfunc(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("")) //nolint:errcheck // ignoring error
}

func weatherHandlerFunc(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	c := q.Get("city")
	we, err := weather.GetWeather(r.Context(), c)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.GetLogger(r.Context()).With("error", err).With(slog.String("city", c)).
			Error("error getting external weather")
		_, _ = w.Write([]byte(err.Error())) //nolint:errcheck // ignoring error
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(we) //nolint:errcheck // ignoring error
}
