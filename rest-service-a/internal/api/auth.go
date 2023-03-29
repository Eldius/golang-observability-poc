package api

import (
    "context"
    "errors"
    "github.com/eldius/golang-observability-poc/otel-instrumentation-helper/telemetry"
    "github.com/jmoiron/sqlx"
    "github.com/rs/zerolog/log"
    "net/http"
    "strings"
)

type User struct {
    ID       int64  `db:"id"`
    Name     string `db:"name"`
    Username string `db:"username"`
    ApiKey   string `db:"api_key"`
}

// AuthApiKey implements a simple middleware handler for adding basic http auth to a route.
func AuthApiKey(realm string, db *sqlx.DB) func(next http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            //l := logger.GetLogger(r.Context())

            ctx := r.Context()
            if r.URL.Path != "/health" {
                authHeader := strings.Trim(r.Header.Get("Authorization"), " ")

                if len(authHeader) == 0 {
                    log.Warn().Msgf("empty auth header")
                    w.WriteHeader(http.StatusUnauthorized)
                    telemetry.NotifyError(r.Context(), errors.New("unauthenticated request"))
                    return
                }

                var results []User
                if err := db.SelectContext(
                    r.Context(),
                    &results,
                    "select id, name, username, api_key from api_users where api_key = $1",
                    authHeader,
                ); err != nil {
                    log.Error().Err(err).Str("api_key", authHeader).Msgf("failed to query db")
                    w.WriteHeader(http.StatusInternalServerError)
                    telemetry.NotifyError(r.Context(), err)
                    return
                }
                if len(results) != 1 {
                    log.Warn().Str("api_key", authHeader).Msgf("wrong query results count: %d", len(results))
                    w.WriteHeader(http.StatusUnauthorized)
                    telemetry.NotifyError(r.Context(), errors.New("unauthorized request"))
                    return
                }
                ctx = context.WithValue(r.Context(), "user", results[0])

                log.Debug().Str("api_key", authHeader).Msgf("right query results count: %d", len(results))
            }

            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
