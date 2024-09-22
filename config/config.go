package config

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/joho/godotenv"

	"github.com/karatekaneen/stockybot/bot"
	"github.com/karatekaneen/stockybot/db"
	"github.com/karatekaneen/stockybot/firestore"
	"github.com/karatekaneen/stockybot/predictor"
)

type Config struct {
	Predictor predictor.Config `embed:""`
	FireDB    firestore.Config `embed:""`
	SQLDB     db.Config        `embed:""`
	Log       LogConfig        `embed:""`
	Bot       bot.Config       `embed:""`
}

//nolint:revive
type LogConfig struct {
	Env string `help:"Log verbosity, set to 'production' if you want json and info" default:"dev" env:"LOG_ENV"`
}

func Load(envPaths ...string) (Config, error) {
	for _, p := range envPaths {
		godotenv.Load(p) //nolint:errcheck
	}

	c := new(Config)

	if err := kong.Parse(c).Validate(); err != nil {
		return *c, fmt.Errorf("config validation: %w", err)
	}

	return *c, nil
}
