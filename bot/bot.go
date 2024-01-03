// Package bot handles all the interactions with Discord
package bot

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/karatekaneen/stockybot/predictor"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type DiscordBot struct {
	predictor         *predictor.Predictor
	dataRepository    dataRepository
	watchRepo         subscriptionRepository
	session           *discordgo.Session
	log               *zap.SugaredLogger
	defaultStockLists map[string]struct{}
	commands          []*discordgo.ApplicationCommand
	cfg               Config
}

type Config struct {
	DefaultTimeout time.Duration `help:"Default timeout for operations"            default:"60s"   env:"DEFAULT_TIMEOUT"`
	RemoveCommands bool          `help:"If commands should be removed on shutdown" default:"true"  env:"REMOVE_COMMANDS"`
	IndexId        int64         `help:"The ID of the index to use as benchmark"   default:"19002" env:"MARKET_INDEX_ID"`
	Token          string        `help:"Auth token"                                                env:"TOKEN"           required:""`
	GuildID        string        `help:"Guild ID to connect to"                                    env:"GUILD_ID"        required:""`
}

func NewBot(
	config Config,
	log *zap.SugaredLogger,
	repo dataRepository,
	pred *predictor.Predictor,
) (*DiscordBot, error) {
	session, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		return nil, errors.Wrap(err, "Invalid bot parameters")
	}

	bot := &DiscordBot{
		session:        session,
		cfg:            config,
		log:            log,
		dataRepository: repo,
		predictor:      pred,
		defaultStockLists: map[string]struct{}{
			"Small Cap Stockholm":  {},
			"Mid Cap Stockholm":    {},
			"Large Cap Stockholm":  {},
			"Large Cap Copenhagen": {},
		},
	}

	// Register handlers
	bot.registerHandlers()

	// Authenticate
	if err := bot.authenticate(); err != nil {
		return nil, errors.Wrap(err, "Could not authenticate")
	}

	// Register commands
	if err := bot.registerCommands(bot.listCommands()); err != nil {
		return nil, errors.Wrap(err, "Could not authenticate")
	}

	return bot, nil
}

// registerHandlers adds functionality similar to a router where it maps
// the incoming command to its designated handler
func (bot *DiscordBot) registerHandlers() {
	handlers := bot.getHandlers()

	bot.session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if handleFunc, ok := handlers[i.ApplicationCommandData().Name]; ok {
			handleFunc(s, i) // TODO: Add errors here
		}
	})
}

// registerCommands lets the Discord server know what functionality the bot provides
func (bot *DiscordBot) registerCommands(commands []*discordgo.ApplicationCommand) error {
	bot.log.Info("Adding commands...")

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))

	for i, rawCmd := range commands {
		cmd, err := bot.session.ApplicationCommandCreate(
			bot.session.State.User.ID,
			bot.cfg.GuildID,
			rawCmd,
		)
		if err != nil {
			return errors.Wrapf(err, "Cannot create '%v'", rawCmd.Name)
		}

		registeredCommands[i] = cmd
	}

	bot.commands = registeredCommands

	return nil
}

func (bot *DiscordBot) authenticate() error {
	bot.session.AddHandler(func(s *discordgo.Session, _ *discordgo.Ready) {
		bot.log.Infof("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	if err := bot.session.Open(); err != nil {
		return errors.Wrap(err, "could not open session")
	}

	return nil
}

func (bot *DiscordBot) Dispose() error {
	defer bot.session.Close()

	if bot.cfg.RemoveCommands {
		bot.log.Info("Removing commands...")

		for _, cmd := range bot.commands {
			err := bot.session.ApplicationCommandDelete(
				bot.session.State.User.ID,
				bot.cfg.GuildID,
				cmd.ID,
			)
			if err != nil {
				return errors.Wrapf(err, "Cannot delete '%v'", cmd.Name)
			}
		}
	}

	return nil
}
