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

	"github.com/qury53/GC_learn/agent/internal/app"
	"github.com/qury53/GC_learn/platform/logger"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatal("failed to load env", err)
	}

	//config := ai_provider_gc.NewConfigMust()
	loggerConfig := logger.NewConfigMust()
	logger := logger.New(loggerConfig, "agent")

	appConfig := app.NewConfigMust()
	app, err := app.New(ctx, logger, appConfig)
	if err != nil {
		logger.Error(
			"failed to create agent service",
			slog.Any("error", err),
		)

		os.Exit(1)
	}

	if err := app.Run(ctx); err != nil {
		logger.Error(
			"failed to run",
			slog.Any("error", err),
		)

		os.Exit(1)
	}
}
