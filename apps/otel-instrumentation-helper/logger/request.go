package logger

import (
	log "github.com/chi-middleware/logrus-logger"
	"github.com/go-chi/chi/v5"
)

func SetupRequestLog(r *chi.Mux) {
	l := Logger()

	r.Use(log.Logger("router", l))
}
