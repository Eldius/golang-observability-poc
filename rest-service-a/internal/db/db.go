package db

import (
	"fmt"
	"github.com/eldius/rest-api/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	"github.com/uptrace/opentelemetry-go-extra/otelsqlx"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
)

var db *sqlx.DB

func DB() *sqlx.DB {
	if db == nil {
		_db, err := open()
		if err != nil {
			log.Fatal().
				Err(err).
				Msg("failed to create db pool")
		}
		db = _db
	}

	return db
}

func open() (*sqlx.DB, error) {

	if config.EnableTraceDB() {
		return otelsqlx.Open(
			"postgres",
			getDBConnectionString(),
			otelsql.WithAttributes(semconv.DBSystemPostgreSQL),
			otelsql.WithDBName(config.GetDBName()),
		)
	}
	return sqlx.Open(
		"postgres",
		getDBConnectionString(),
	)
}

func getDBConnectionString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.GetDBUser(),
		config.GetDBPass(),
		config.GetDBHost(),
		config.GetDBPort(),
		config.GetDBName(),
	)
}
