package main

import (
	"context"
	"log"
	"sort"

	"github.com/karatekaneen/stockybot/config"
	"github.com/karatekaneen/stockybot/db"
	"github.com/karatekaneen/stockybot/ent/security"
	"github.com/lithammer/fuzzysearch/fuzzy"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

	zaplog, err := createLogger(cfg.Log)
	if err != nil {
		log.Fatal(err)
	}

	logger := zaplog.Sugar()

	ctx := context.Background()

	sqlDB, err := db.New(ctx, cfg.SQLDB, logger)
	if err != nil {
		logger.Errorf("failed sql db init: %w", err)
		return
	}

	defer sqlDB.Close()

	names := []string{}

	for _, s := range sqlDB.Client.Security.Query().Select(security.FieldName).AllX(ctx) {
		names = append(names, s.Name)
	}

	logger.Info(names)

	result := fuzzy.RankFindNormalizedFold("invest", names)
	sort.Sort(result)

	for i, item := range result {
		logger.Info(i, item.Target)
	}
}
