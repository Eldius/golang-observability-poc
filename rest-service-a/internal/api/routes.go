package api

import (
	"fmt"
	"github.com/eldius/rest-api/internal/config"
	"github.com/eldius/rest-api/internal/db"
	"github.com/eldius/rest-api/internal/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
	"github.com/riandyrn/otelchi"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
)

func Start(port int) {

	httpLogger := httplog.NewLogger(config.GetServiceName(), httplog.Options{
		JSON: true,
	})

	r := chi.NewRouter()

	setupRestTracing(r)

	r.Use(httplog.RequestLogger(httpLogger))
	r.Use(AuthApiKey("api", db.DB()))

	r.Get("/", homeHandlerfunc)
	r.Get("/ping", pingHandlerfunc)
	r.Get("/health", healthHandlerfunc)

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

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

func pingHandlerfunc(w http.ResponseWriter, r *http.Request) {

	l := logger.GetLogger(r.Context())
	l.Info().Msg("get ping begin")
	l.Info().Msgf("testing: %s", r.Context())

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("welcome"))

	l.Info().Msg("get ping end")
}

func healthHandlerfunc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(""))
}

func setupRestTracing(r *chi.Mux) {
	if config.EnableTraceRest() {
		r.Use(otelchi.Middleware(config.GetServiceName(), otelchi.WithChiRoutes(r)))
	}
}
