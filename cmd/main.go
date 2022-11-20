package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/karatekaneen/stockybot/bot"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Bot parameters
var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	BotToken       = flag.String("token", "", "Bot access token")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

type LogConfig struct {
	Env string
}

func init() { flag.Parse() }

func main() {
	l, _ := createLogger(LogConfig{})
	logger := l.Sugar()

	b, err := bot.NewBot(*BotToken, *GuildID, *RemoveCommands, l.Sugar())
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

func createLogger(cfg LogConfig) (*zap.Logger, error) {
	var loggerSettings zap.Config

	if cfg.Env == "production" {
		loggerSettings = zap.NewProductionConfig()
		loggerSettings.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else {
		loggerSettings = zap.NewDevelopmentConfig()
		loggerSettings.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	return loggerSettings.Build()
}
