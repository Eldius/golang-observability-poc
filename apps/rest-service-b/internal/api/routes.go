package api

import (
	"encoding/json"
	"fmt"
	"github.com/eldius/golang-observability-poc/apps/otel-instrumentation-helper/logger"
	"github.com/eldius/golang-observability-poc/apps/otel-instrumentation-helper/telemetry"
	"github.com/eldius/golang-observability-poc/apps/rest-service-b/internal/config"
	"github.com/eldius/golang-observability-poc/apps/rest-service-b/internal/weather"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
	"net/http"
	"time"
)

func Start(port int) {

	l := logger.Logger()

	httpLogger := httplog.NewLogger(config.GetServiceName(), httplog.Options{
		JSON: true,
	})

	r := chi.NewRouter()

	telemetry.SetupRestTracing(r)

	r.Use(httplog.RequestLogger(httpLogger))

	r.Get("/", homeHandlerfunc)
	r.Get("/health", healthHandlerfunc)
	r.Get("/weather", weatherHandlerFunc)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           r,
		ReadHeaderTimeout: 100 * time.Millisecond,
	}
	l.Info().
		Str("addr", srv.Addr).
		Msg("starting")
	if err := srv.ListenAndServe(); err != nil {
		l.Fatal().Err(err).Msg("filed to start server")
	}
}

func homeHandlerfunc(w http.ResponseWriter, r *http.Request) {

	l := logger.GetLogger(r.Context())
	l.Info().Msg("get root begin")
	l.Info().Msgf("testing: %s", r.Context())

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("home")) //nolint:errcheck // ignoring error

	l.Info().Msg("get root end")
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
		_, _ = w.Write([]byte(err.Error())) //nolint:errcheck // ignoring error
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(we) //nolint:errcheck // ignoring error
}
