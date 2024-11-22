package main

import (
	"context"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/karatekaneen/stockybot/config"
	"github.com/karatekaneen/stockybot/firestore"
)

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

func main() {
	cfg, err := config.Load(".env", ".env.local", ".local.env")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	db, err := firestore.New(ctx, cfg.FireDB)
	if err != nil {
		panic(err)
	}

	dailyCtx, err := db.StrategyState(ctx, 5265)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n\n", dailyCtx)
}
