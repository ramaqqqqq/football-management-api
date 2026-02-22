package main

import (
	"context"
	"os"

	"go-test/lib/logger"
	"go-test/migration"
	"go-test/src/app"
)

func main() {
	ctx := context.Background()

	logger.Init(ctx)

	if err := app.Init(ctx); err != nil {
		logger.GetLogger(ctx).Fatalf("Failed to initialize app: %v", err)
	}
	defer app.Close()

	args := os.Args
	if len(args) < 2 {
		logger.GetLogger(ctx).Fatal("Missing args. args: [up | rollback]")
	}

	migrationSvc, err := migration.New(ctx, app.Config().Postgres)
	if err != nil {
		logger.GetLogger(ctx).Fatalf("Failed to initiate migration: %v", err)
	}

	switch args[1] {
	case "up":
		if err := migrationSvc.Up(ctx); err != nil {
			logger.GetLogger(ctx).Fatalf("Failed to run migration up: %v", err)
		}
		logger.GetLogger(ctx).Info("Migration up completed successfully")
	case "rollback":
		if err := migrationSvc.Rollback(ctx); err != nil {
			logger.GetLogger(ctx).Fatalf("Failed to rollback migration: %v", err)
		}
		logger.GetLogger(ctx).Info("Migration rollback completed successfully")
	default:
		logger.GetLogger(ctx).Fatal("Invalid migration command. Use: [up | rollback]")
	}
}
