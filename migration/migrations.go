package migration

import (
	"context"
	"errors"

	"go-test/lib/logger"
	"go-test/src/app"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	pgMigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
)

const (
	migrateLogIdentifier = "go-test"
)

type MigrationService interface {
	Up(context.Context) error
	Rollback(context.Context) error
	Version(context.Context) (int, bool, error)
}

type migrationService struct {
	driver  database.Driver
	migrate *migrate.Migrate
}

func New(ctx context.Context, cfg app.Postgres) (MigrationService, error) {
	db, err := sqlx.Connect("postgres", cfg.ConnURI)
	if err != nil {
		logger.GetLogger(ctx).Errorf("error connecting to postgres url:%s, err: %v", cfg.ConnURI, err)
		return nil, err
	}

	databaseInstance, err := pgMigrate.WithInstance(db.DB, &pgMigrate.Config{})
	if err != nil {
		logger.GetLogger(ctx).Errorf("go-migrate postgres drv init failed: %v", err)
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance("file://migration/sql",
		migrateLogIdentifier, databaseInstance)
	if err != nil {
		logger.GetLogger(ctx).Errorf("migrate init failed %v", err)
		return nil, err
	}

	return migrationService{
		driver:  databaseInstance,
		migrate: m,
	}, nil
}

func (s migrationService) Up(ctx context.Context) error {
	currVersion, _, err := s.Version(ctx)
	if err != nil {
		logger.GetLogger(ctx).Error("Failed get current version err: ", err)
		return err
	}

	logger.GetLogger(ctx).Infof("Running migration from version: %d", currVersion)
	if err := s.migrate.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.GetLogger(ctx).Info("No Changes")
			return nil
		}
		logger.GetLogger(ctx).Error("Failed run migrate err: ", err)
		return err
	}

	currVersion, _, _ = s.Version(ctx)
	logger.GetLogger(ctx).Info("Migration success, current version:", currVersion)
	return nil
}

func (s migrationService) Rollback(ctx context.Context) error {
	currVersion, _, err := s.Version(ctx)
	if err != nil {
		logger.GetLogger(ctx).Error("Failed get current version err: ", err)
		return err
	}

	if currVersion < 0 {
		logger.GetLogger(ctx).Info("No migrations to rollback")
		return nil
	}

	logger.GetLogger(ctx).Infof("Rollingback 1 step from version: %d", currVersion)

	if err := s.migrate.Steps(-1); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.GetLogger(ctx).Info("No Changes")
			return nil
		}
		logger.GetLogger(ctx).Errorf("Failed to rollback, err:%v", err)
		return err
	}

	currVersion, _, _ = s.Version(ctx)
	logger.GetLogger(ctx).Infof("Rollback success, current version:%d", currVersion)
	return nil
}

func (s migrationService) Version(ctx context.Context) (int, bool, error) {
	currVersion, dirty, err := s.driver.Version()
	if err != nil {
		logger.GetLogger(ctx).Errorf("Failed to get version: %v", err)
		return 0, false, err
	}
	return currVersion, dirty, nil
}
