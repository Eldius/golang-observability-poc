package db

import (
    "context"
    "github.com/eldius/golang-observability-poc/rest-service-a/internal/config"
    _ "github.com/lib/pq"
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
    migrate "github.com/rubenv/sql-migrate"
)

func Migrations() error {
    if config.GetMigrationsEnabled() {
        db := DB()

        migrations := &migrate.FileMigrationSource{
            Dir: "db/migrations",
        }

        migInfo, err := migrations.FindMigrations()
        log.Info().Int("migrations_to_do", len(migInfo)).Msgf("running migrations begin")

        n, err := migrate.Exec(db.DB, "postgres", migrations, migrate.Up)
        if err != nil {
            // Handle errors!
            log.Error().Err(err).Msg("failed to run migrations")
        }
        log.Info().Int("migrations_done", n).Msgf("running migrations end")
        //fmt.Printf("Applied %d migrations!\n", n)

        //driver, err := postgres.WithInstance(db.DB, &postgres.Config{
        //    MigrationsTable: postgres.DefaultMigrationsTable,
        //    DatabaseName:    config.GetDBName(),
        //    SchemaName:      "public",
        //})
        //if err != nil {
        //    return err
        //}
        //
        //m, err := migrate.NewWithDatabaseInstance(
        //    migrationsScriptsPath,
        //    "postgres", driver)
        //if err != nil {
        //    return err
        //}
        //
        //if isMigrationLogEnabled() {
        //    m.Log = newMigrationLogger()
        //}
        //
        //if err := m.Up(); err != nil {
        //    log.Error().Err(err).Msg("failed to run migrations")
        //}
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
