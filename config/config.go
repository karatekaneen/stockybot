package config

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/joho/godotenv"
	"github.com/karatekaneen/stockybot/bot"
	"github.com/karatekaneen/stockybot/firestore"
	"github.com/karatekaneen/stockybot/predictor"
)

type Config struct {
	Bot       bot.Config       `embed:""`
	Predictor predictor.Config `embed:""`
	DB        firestore.Config `embed:""`
	Log       LogConfig        `embed:""`
}

type LogConfig struct {
	Env string `help:"Log verbosity, set to 'production' if you want json and info" default:"dev" env:"LOG_ENV"`
}

func Load(envPaths ...string) (Config, error) {
	for _, p := range envPaths {
		godotenv.Load(p)
	}

	c := Config{}

	if err := kong.Parse(&c).Validate(); err != nil {
		return c, fmt.Errorf("config validation: %w", err)
	}

	return c, nil
}
