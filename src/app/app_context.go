package app

import (
	"context"
	"go-test/lib/i18n"
	"go-test/lib/logger"
	"time"

	"github.com/go-playground/validator/v10"
	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
)

type appContext struct {
	db               *sqlx.DB
	requestValidator *validator.Validate
	cfg              *Configuration
}

var appCtx appContext

func Init(ctx context.Context) error {
	logger.Init(ctx)

	cfg, err := InitConfig(ctx)
	if err != nil {
		return err
	}

	if err := i18n.Init(ctx, cfg.Translation.FilePath, cfg.Translation.DefaultLanguage); err != nil {
		panic(err)
	}

	db, err := sqlx.Connect("postgres", cfg.Postgres.ConnURI)
	if err != nil {
		return err
	}

	if cfg.Postgres.MaxPoolSize > 0 {
		db.SetMaxOpenConns(cfg.Postgres.MaxPoolSize)
	}
	if cfg.Postgres.MaxIdleConnections > 0 {
		db.SetMaxIdleConns(cfg.Postgres.MaxIdleConnections)
	}
	if cfg.Postgres.MaxIdleTime > 0 {
		db.SetConnMaxIdleTime(cfg.Postgres.MaxIdleTime * time.Second)
	}
	if cfg.Postgres.MaxLifeTime > 0 {
		db.SetConnMaxLifetime(cfg.Postgres.MaxLifeTime * time.Second)
	}

	appCtx = appContext{
		db:               db,
		requestValidator: validator.New(),
		cfg:              cfg,
	}

	return nil
}

func RequestValidator() *validator.Validate {
	return appCtx.requestValidator
}

func DB() *sqlx.DB {
	return appCtx.db
}

func Config() Configuration {
	return *appCtx.cfg
}

func Close() error {
	if appCtx.db != nil {
		return appCtx.db.Close()
	}
	return nil
}
