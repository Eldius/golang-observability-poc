package db

import (
	"github.com/eldius/golang-observability-poc/otel-instrumentation-helper/logger"
	"github.com/eldius/golang-observability-poc/rest-service-a/internal/config"
	_ "github.com/lib/pq" // we need the Postgres driver
	"github.com/pkg/errors"
	migrate "github.com/rubenv/sql-migrate"
	"log/slog"
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
			err = errors.Wrap(err, "failed to find migrations")
			l.With("error", err).Error("failed to find migrations")
			return err
		}
		l.With(slog.Int("migrations_to_do", len(migInfo))).Info("running migrations begin")

		n, err := migrate.Exec(db.DB, "postgres", migrations, migrate.Up)
		if err != nil {
			err = errors.Wrap(err, "failed to execute migrations")
			l.With("error", err).Error("failed to run migrations")
			return err
		}
		l.With(slog.Int("migrations_done", n)).Info("running migrations end")
	}

	return nil
}
