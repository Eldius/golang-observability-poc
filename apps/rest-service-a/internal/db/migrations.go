package db

import (
	"github.com/eldius/golang-observability-poc/apps/otel-instrumentation-helper/logger"
	"github.com/eldius/golang-observability-poc/apps/rest-service-a/internal/config"
	_ "github.com/lib/pq" // we need the Postgres driver
	"github.com/pkg/errors"
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
			err = errors.Wrap(err, "failed to find migrations")
			l.WithError(err).Error("failed to find migrations")
			return err
		}
		l.WithField("migrations_to_do", len(migInfo)).Infof("running migrations begin")

		n, err := migrate.Exec(db.DB, "postgres", migrations, migrate.Up)
		if err != nil {
			err = errors.Wrap(err, "failed to execute migrations")
			l.WithError(err).Error("failed to run migrations")
			return err
		}
		l.WithField("migrations_done", n).Info("running migrations end")
	}

	return nil
}
