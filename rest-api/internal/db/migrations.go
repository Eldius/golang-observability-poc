package db

import (
	"context"
	"fmt"
	"github.com/eldius/rest-api/internal/config"
	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func Migrations() error {
	if config.GetMigrationsEnabled() {
		wd, _ := os.Getwd()
		migrationsScriptsPath := fmt.Sprintf("file://%s/%s", wd, "db/migrations/")
		log.Info().Str("migrations_path", migrationsScriptsPath).Msgf("running migrations")

		db := DB()

		driver, err := postgres.WithInstance(db.DB, &postgres.Config{
			MigrationsTable: postgres.DefaultMigrationsTable,
			DatabaseName:    config.GetDBName(),
			SchemaName:      "public",
		})
		if err != nil {
			return err
		}

		m, err := migrate.NewWithDatabaseInstance(
			migrationsScriptsPath,
			"postgres", driver)
		if err != nil {
			return err
		}

		if isMigrationLogEnabled() {
			m.Log = newMigrationLogger()
		}

		if err := m.Up(); err != nil {
			log.Error().Err(err).Msg("failed to run migrations")
		}
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
	return &migrationLogger{
		l: zerolog.Ctx(context.Background()),
		v: true,
	}
}

func isMigrationLogEnabled() bool {
	return zerolog.GlobalLevel() >= zerolog.DebugLevel
}
