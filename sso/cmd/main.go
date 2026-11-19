package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/qury53/GC_learn/platform/logger"
	application "github.com/qury53/GC_learn/sso/internal/app"
)

const (
	serviceName = "sso"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	loggerConfig := logger.NewConfigMust()
	log := logger.New(loggerConfig, serviceName)

	appConfig := application.NewConfigMust()
	app, err := application.New(ctx, appConfig, log)

	if err != nil {
		log.Error(
			"application initializing failed",
			slog.Any("error", err),
		)

		os.Exit(1)
	}

	if err := app.Run(ctx); err != nil {
		log.Error(
			"application run error",
			slog.Any("error", err),
		)

		os.Exit(1)
	}
}
