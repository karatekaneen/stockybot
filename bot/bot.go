// Package bot handles all the interactions with Discord
package bot

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"

	"github.com/karatekaneen/stockybot"
	"github.com/karatekaneen/stockybot/predictor"
)

type subscriptionRepository interface {
	AddSubscription(ctx context.Context, stockName string, userID string) error
	RemoveSubscription(ctx context.Context, stockName string, userID string) error
	GetSubscribedSecurities(ctx context.Context, userID string) ([]stockybot.Security, error)
	GetAllStockNames(ctx context.Context) ([]string, error)
}

type DiscordBot struct {
	session  *discordgo.Session
	log      *zap.SugaredLogger
	router   *router
	commands []*discordgo.ApplicationCommand
	cfg      Config
}

//nolint:revive
type Config struct {
	Token          string        `help:"Auth token"                                env:"TOKEN"           required:""`
	GuildID        string        `help:"Guild ID to connect to"                    env:"GUILD_ID"        required:""`
	DefaultTimeout time.Duration `help:"Default timeout for operations"            env:"DEFAULT_TIMEOUT"             default:"60s"`
	IndexID        int64         `help:"The ID of the index to use as benchmark"   env:"MARKET_INDEX_ID"             default:"19002"`
	RemoveCommands bool          `help:"If commands should be removed on shutdown" env:"REMOVE_COMMANDS"             default:"true"`
}

func NewBot(
	config Config,
	log *zap.SugaredLogger,
	repo dataRepository,
	pred *predictor.Predictor,
	watchRepo subscriptionRepository,
) (*DiscordBot, error) {
	session, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		return nil, fmt.Errorf("instantiate bot: %w", err)
	}

	router := newRouter(config, log, repo, pred, watchRepo)

	bot := &DiscordBot{
		session: session,
		cfg:     config,
		log:     log,
		router:  router,
	}

	// Register handlers
	bot.session.AddHandler(router.Handle)

	// Authenticate
	if err := bot.authenticate(); err != nil {
		return nil, fmt.Errorf("discord authentication: %w", err)
	}

	// Register commands
	if err := bot.registerCommands(router.Commands()); err != nil {
		return nil, fmt.Errorf("register commands: %w", err)
	}

	_, err = bot.session.ChannelMessageSend(
		"1043847870861811807",
		"hello there <@344808325394726912>",
	)
	if err != nil {
		return nil, err
	}

	return bot, nil
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
			return fmt.Errorf("create command %q: %w", rawCmd.Name, err)
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
		return fmt.Errorf("open session: %w", err)
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
				return fmt.Errorf("delete command %q: %w", cmd.Name, err)
			}
		}
	}

	return nil
}
