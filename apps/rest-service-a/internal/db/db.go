package db

import (
	"fmt"
	"github.com/eldius/golang-observability-poc/apps/otel-instrumentation-helper/logger"
	"github.com/eldius/golang-observability-poc/apps/rest-service-a/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // we need the Postgres driver
	"github.com/pkg/errors"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	"github.com/uptrace/opentelemetry-go-extra/otelsqlx"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
)

var db *sqlx.DB

func DB() *sqlx.DB {
	if db == nil {
		_db, err := open()
		if err != nil {
			err = errors.Wrap(err, "failed to open db connection")
			logger.Logger().
				WithError(err).
				Fatal("failed to create db pool")
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
