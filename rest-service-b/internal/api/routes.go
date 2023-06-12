package api

import (
	"encoding/json"
	"fmt"
	"github.com/eldius/golang-observability-poc/otel-instrumentation-helper/logger"
	"github.com/eldius/golang-observability-poc/otel-instrumentation-helper/telemetry"
	"github.com/eldius/golang-observability-poc/rest-service-b/internal/weather"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func Start(port int) {

	l := logger.Logger()

	r := chi.NewRouter()

	telemetry.SetupRestTracing(r)
	logger.SetupRequestLog(r)

	r.Get("/", homeHandlerfunc)
	r.Get("/health", healthHandlerfunc)
	r.Get("/weather", weatherHandlerFunc)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           r,
		ReadHeaderTimeout: 100 * time.Millisecond,
	}
	l.WithField("addr", srv.Addr).
		Info("starting")
	if err := srv.ListenAndServe(); err != nil {
		l.WithError(err).Fatal("filed to start server")
	}
}

func homeHandlerfunc(w http.ResponseWriter, r *http.Request) {

	l := logger.GetLogger(r.Context())
	l.Info("get root begin")
	l.Infof("testing: %s", r.Context())

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
		logger.GetLogger(r.Context()).WithFields(logrus.Fields{
			"city": c,
		}).WithError(err).
			Error("error getting external weather")
		_, _ = w.Write([]byte(err.Error())) //nolint:errcheck // ignoring error
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(we) //nolint:errcheck // ignoring error
}
