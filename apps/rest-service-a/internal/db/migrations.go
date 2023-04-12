package db

import (
	"github.com/eldius/golang-observability-poc/apps/otel-instrumentation-helper/logger"
	"github.com/eldius/golang-observability-poc/apps/rest-service-a/internal/config"
	_ "github.com/lib/pq" // we need the Postgres driver
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
		if err != nil {
			l.Error().Err(err).Msg("failed to find migrations")
			return err
		}
		l.Info().Int("migrations_to_do", len(migInfo)).Msgf("running migrations begin")

		n, err := migrate.Exec(db.DB, "postgres", migrations, migrate.Up)
		if err != nil {
			l.Error().Err(err).Msg("failed to run migrations")
			return err
		}
		l.Info().Int("migrations_done", n).Msgf("running migrations end")
	}

	return nil
}
