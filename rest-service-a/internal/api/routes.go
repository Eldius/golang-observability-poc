package api

import (
	"encoding/json"
	"fmt"
	"github.com/eldius/golang-observability-poc/otel-instrumentation-helper/logger"
	"github.com/eldius/golang-observability-poc/otel-instrumentation-helper/telemetry"
	"github.com/eldius/golang-observability-poc/rest-service-a/internal/db"
	"github.com/eldius/golang-observability-poc/rest-service-a/internal/integration/serviceb"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
)

func Start(port int) {

	l := logger.Logger()

	r := chi.NewRouter()

	telemetry.SetupRestTracing(r)
	logger.SetupRequestLog(r)

	r.Use(AuthAPIKey("api", db.DB()))

	r.Get("/", homeHandlerfunc)
	r.Get("/ping", pingHandlerfunc)
	r.Get("/health", healthHandlerfunc)
	r.Get("/weather", weatherHandlerfunc)

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

func pingHandlerfunc(w http.ResponseWriter, r *http.Request) {

	l := logger.GetLogger(r.Context())
	l.Info("get ping begin")
	l.Infof("testing: %s", r.Context())

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("welcome")) //nolint:errcheck // ignoring error

	l.Info("get ping end")
}

func healthHandlerfunc(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("")) //nolint:errcheck // ignoring error
}

func weatherHandlerfunc(w http.ResponseWriter, r *http.Request) {
	l := logger.GetLogger(r.Context())
	q := r.URL.Query()
	c := q.Get("city")
	if c == "" {
		l.Error("city is empty")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("city is empty")) //nolint:errcheck // ignoring error
		return
	}
	we, err := serviceb.GetWeather(r.Context(), c)
	if err != nil {
		l.WithError(err).WithField("city", c).Error("service-b integration failed")
		telemetry.NotifyError(r.Context(), err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error())) //nolint:errcheck // ignoring error
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(&we) //nolint:errcheck // ignoring error
}
