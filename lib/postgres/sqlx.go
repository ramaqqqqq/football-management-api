package postgres

import (
	"context"
	"go-test/lib/logger"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	"github.com/uptrace/opentelemetry-go-extra/otelsqlx"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func InitSQLX(ctx context.Context, cfg PostgresConfig) (*sqlx.DB, error) {
	xDB, err := otelsqlx.Open("postgres", cfg.ConnectionUrl, otelsql.WithAttributes(semconv.DBSystemPostgreSQL))
	if err != nil {
		logger.GetLogger(ctx).Errorf("failed to load the database err:%v", err)
		return nil, err
	}

	if err = xDB.Ping(); err != nil {
		logger.GetLogger(ctx).Errorf("failed to ping the database err:%v", err)
		return nil, err
	}

	xDB.SetMaxOpenConns(cfg.MaxPoolSize)
	xDB.SetMaxIdleConns(cfg.MaxIdleConnections)
	xDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	xDB.SetConnMaxLifetime(cfg.ConnMaxLifeTime)
	return xDB, nil
}
