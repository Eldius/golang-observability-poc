package api

import (
	"context"
	"errors"
	"github.com/eldius/golang-observability-poc/apps/otel-instrumentation-helper/logger"
	"github.com/eldius/golang-observability-poc/apps/otel-instrumentation-helper/telemetry"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
)

type User struct {
	ID       int64  `db:"id"`
	Name     string `db:"name"`
	Username string `db:"username"`
	APIKey   string `db:"api_key"`
}

type RequestContextKey struct {
	name string
}

var UserContextKey = RequestContextKey{name: "user"}

// AuthAPIKey implements a simple middleware handler for adding basic http auth to a route.
func AuthAPIKey(_ string, db *sqlx.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l := logger.GetLogger(r.Context())

			ctx := r.Context()

			ctx, f := telemetry.StartSpan(ctx, "UserValidation")
			defer f()

			if r.URL.Path != "/health" {
				authHeader := strings.Trim(r.Header.Get("Authorization"), " ")

				if authHeader == "" {
					log.Warn().Msgf("empty auth header")
					w.WriteHeader(http.StatusUnauthorized)
					telemetry.NotifyError(r.Context(), errors.New("unauthenticated request"))
					return
				}

				var results []User
				if err := db.SelectContext(
					ctx,
					&results,
					"select id, name, username, api_key from api_users where api_key = $1",
					authHeader,
				); err != nil {
					l.Error().Err(err).Str("api_key", authHeader).Msgf("failed to query db")
					w.WriteHeader(http.StatusInternalServerError)
					telemetry.NotifyError(r.Context(), err)
					return
				}
				if len(results) != 1 {
					l.Warn().Str("api_key", authHeader).Msgf("wrong query results count: %d", len(results))
					w.WriteHeader(http.StatusUnauthorized)
					telemetry.NotifyError(r.Context(), errors.New("unauthorized request"))
					return
				}
				ctx = context.WithValue(r.Context(), UserContextKey, results[0])

				telemetry.AddAttribute(ctx, "user", results[0].Username)

				l.Debug().Str("api_key", authHeader).Msgf("right query results count: %d", len(results))
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
