package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/karatekaneen/stockybot/bot"
	"github.com/karatekaneen/stockybot/config"
	"github.com/karatekaneen/stockybot/firestore"
	"github.com/karatekaneen/stockybot/predictor"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	cfg, err := config.Load(".env", ".env.local", ".local.env")
	if err != nil {
		log.Fatal(err)
	}

	l, err := createLogger(cfg.Log)
	if err != nil {
		log.Fatal(err)
	}
	logger := l.Sugar()

	ctx := context.Background()

	db, err := firestore.New(ctx, cfg.DB)
	if err != nil {
		logger.Fatalln(err)
	}

	b, err := bot.NewBot(cfg.Bot, l.Sugar(), db, predictor.New(cfg.Predictor, logger))
	if err != nil {
		logger.Fatal(err)
	}

	defer b.Dispose()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	logger.Info("Press Ctrl+C to exit")
	<-stop

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

	return loggerSettings.Build()
}
