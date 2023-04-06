package db

import (
	"context"
	"github.com/eldius/golang-observability-poc/apps/otel-instrumentation-helper/logger"
	"github.com/eldius/golang-observability-poc/apps/rest-service-a/internal/config"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	migrate "github.com/rubenv/sql-migrate"
)

func Migrations() error {
	if config.GetMigrationsEnabled() {
		l := logger.Logger()
		db := DB()

		migrations := &migrate.FileMigrationSource{
			Dir: "db/migrations",
		}

		migInfo, err := migrations.FindMigrations()
		l.Info().Int("migrations_to_do", len(migInfo)).Msgf("running migrations begin")

		n, err := migrate.Exec(db.DB, "postgres", migrations, migrate.Up)
		if err != nil {
			// Handle errors!
			l.Error().Err(err).Msg("failed to run migrations")
		}
		l.Info().Int("migrations_done", n).Msgf("running migrations end")
	}

	return nil
}

type migrationLogger struct {
	l *zerolog.Logger
	v bool
}

func (m *migrationLogger) Printf(format string, v ...interface{}) {
	m.l.Debug().Msgf(format, v)
}
func (m *migrationLogger) Verbose() bool {
	return m.v
}

func newMigrationLogger() *migrationLogger {
	l := logger.GetLogger(context.Background())
	return &migrationLogger{
		l: &l,
		v: true,
	}
}

func isMigrationLogEnabled() bool {
	return zerolog.GlobalLevel() >= zerolog.DebugLevel
}
