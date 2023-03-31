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
	"github.com/rs/zerolog/log"
	"net/http"
)

func Start(port int) {

	httpLogger := httplog.NewLogger(config.GetServiceName(), httplog.Options{
		JSON: true,
	})

	r := chi.NewRouter()

	telemetry.SetupRestTracing(r)

	r.Use(httplog.RequestLogger(httpLogger))

	r.Get("/", homeHandlerfunc)
	r.Get("/health", healthHandlerfunc)
	r.Get("/weather", weatherHandlerFunc)

	log.Info().
		Int("port", port).
		Msg("starting")
	panic(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}

func homeHandlerfunc(w http.ResponseWriter, r *http.Request) {

	l := logger.GetLogger(r.Context())
	l.Info().Msg("get root begin")
	l.Info().Msgf("testing: %s", r.Context())

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("home"))

	l.Info().Msg("get root end")
}

func healthHandlerfunc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(""))
}

func weatherHandlerFunc(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query()
	c := q.Get("city")
	we, err := weather.GetWeather(r.Context(), c)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(we)
}
