package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/karatekaneen/stockybot/bot"
	"github.com/karatekaneen/stockybot/config"
	"github.com/karatekaneen/stockybot/db"
	"github.com/karatekaneen/stockybot/firestore"
	"github.com/karatekaneen/stockybot/predictor"
)

func main() {
	cfg, err := config.Load(".env", ".env.local", ".local.env")
	if err != nil {
		log.Fatal(err)
	}

	zaplog, err := createLogger(cfg.Log)
	if err != nil {
		log.Fatal(err)
	}

	logger := zaplog.Sugar()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sqlDB, err := db.New(ctx, cfg.SQLDB, logger)
	if err != nil {
		logger.Errorf("failed sql db init: %w", err)
		return
	}

	defer sqlDB.Close()

	fireDB, err := firestore.New(ctx, cfg.FireDB)
	if err != nil {
		logger.Error(err)
		return
	}

	b, err := bot.NewBot(
		cfg.Bot,
		zaplog.Sugar(),
		fireDB,
		predictor.New(cfg.Predictor, logger),
		sqlDB,
	)
	if err != nil {
		logger.Error(err)
		return
	}

	defer b.Dispose() //nolint:errcheck

	errCh := make(chan error, 1)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		errCh <- sqlDB.ImportPeriodically(ctx, fireDB, 24*time.Hour) //nolint:mnd
	}()

	logger.Info("Press Ctrl+C to exit")

	select {
	case <-stop:
	case err := <-errCh:
		logger.Error(err)
	}

	logger.Info("Gracefully shutting down.")
}

func createLogger(cfg config.LogConfig) (*zap.Logger, error) {
	var loggerSettings zap.Config

	if cfg.Env == "production" {
		loggerSettings = zap.NewProductionConfig()
		loggerSettings.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	} else {
		loggerSettings = zap.NewDevelopmentConfig()
		loggerSettings.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	return loggerSettings.Build() //nolint:wrapcheck
}
