package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/yandzee/wait-action/internal/app"
	"github.com/yandzee/wait-action/pkg/config"
)

func main() {
	cfg, err := config.ParseEnv()
	if err != nil {
		panic("failed to initialize application config: " + err.Error())
	}

	logger := initLogger(cfg.IsDebugEnabled)
	ctx := context.Background()

	if err := app.Run(ctx, logger, cfg); err != nil {
		logger.Error("application run failed", "err", err.Error())
		os.Exit(1)
	}
}

func initLogger(isDebug bool) *slog.Logger {
	slogLvl := slog.LevelInfo
	if isDebug {
		slogLvl = slog.LevelDebug
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slogLvl,
	}))

	return logger
}
