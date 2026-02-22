// @title			Football Management API
// @version			1.0
// @description		REST API for managing football teams, players, and matches
// @host			localhost:8080
// @BasePath		/
//
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description					Type "Bearer" followed by a space and the JWT token.
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go-test/lib/logger"
	ginmiddleware "go-test/lib/middleware/gin"
	"go-test/src/app"
	v1 "go-test/src/v1"
)

func main() {
	if err := os.Setenv("TZ", "Asia/Jakarta"); err != nil {
		panic(err)
	}

	initCtx := context.Background()
	logger.Init(initCtx)

	if err := app.Init(initCtx); err != nil {
		panic(err)
	}

	startService(initCtx)
}

func startService(ctx context.Context) {
	defer app.Close()

	address := fmt.Sprintf(":%d", app.Config().BindAddress)
	logger.GetLogger(ctx).Infof("Starting service on %s", address)

	gin.SetMode(func() string {
		if app.Config().Environment == "production" {
			return gin.ReleaseMode
		}
		return gin.DebugMode
	}())

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(ginmiddleware.RequestIDMiddleware())

	deps := v1.Dependencies(ctx)
	v1.Router(r, deps)

	server := &http.Server{
		Addr:         address,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.GetLogger(ctx).Infof("Server ready at %s", address)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.GetLogger(ctx).Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.GetLogger(ctx).Info("Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.GetLogger(ctx).Errorf("Server forced to shutdown: %v", err)
	}

	logger.GetLogger(ctx).Info("Server exited gracefully")
}
