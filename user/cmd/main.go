package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/qury53/GC_learn/platform/logger"
	"github.com/qury53/GC_learn/user/internal/app"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatal("failed to load env", err)
	}

	loggerConfig := logger.NewConfigMust()
	logger := logger.New(loggerConfig, "user")

	appConfig := app.NewConfigMust()
	app, err := app.New(ctx, appConfig, logger)
	if err != nil {
		logger.Error(
			"failed to create user app service",
			slog.Any("error", err),
		)

		os.Exit(1)
	}

	if err := app.Run(ctx); err != nil {
		logger.Error(
			"failed to run user service",
			slog.Any("error", err),
		)

		os.Exit(1)
	}
}
