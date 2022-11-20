// Package bot handles all the interactions with Discord
package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type DiscordBot struct {
	commands       []*discordgo.ApplicationCommand
	session        *discordgo.Session
	handlers       map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
	log            *zap.SugaredLogger
	RemoveCommands bool
	GuildID        string
}

func NewBot(token string, guildId string, removeCommands bool, log *zap.SugaredLogger) (*DiscordBot, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, errors.Wrap(err, "Invalid bot parameters")
	}

	bot := &DiscordBot{
		session:        session,
		handlers:       commandHandlers,
		log:            log,
		RemoveCommands: removeCommands,
	}

	// Register handlers
	bot.registerHandlers()

	// Authenticate
	if err := bot.authenticate(); err != nil {
		return nil, errors.Wrap(err, "Could not authenticate")
	}

	// Register commands
	if err := bot.registerCommands(allCommands); err != nil {
		return nil, errors.Wrap(err, "Could not authenticate")
	}

	return bot, nil
}

// registerHandlers adds functionality similar to a router where it maps
// the incoming command to its designated handler
func (bot *DiscordBot) registerHandlers() {
	bot.session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if handleFunc, ok := bot.handlers[i.ApplicationCommandData().Name]; ok {
			handleFunc(s, i) // TODO: Add errors here
		}
	})
}

// registerCommands lets the Discord server know what functionality the bot provides
func (bot *DiscordBot) registerCommands(commands []*discordgo.ApplicationCommand) error {
	bot.log.Info("Adding commands...")

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))

	for i, rawCmd := range commands {
		cmd, err := bot.session.ApplicationCommandCreate(bot.session.State.User.ID, bot.GuildID, rawCmd)
		if err != nil {
			return errors.Wrapf(err, "Cannot create '%v'", rawCmd.Name)
		}

		registeredCommands[i] = cmd
	}

	bot.commands = registeredCommands

	return nil
}

func (bot *DiscordBot) authenticate() error {
	bot.session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		bot.log.Infof("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	if err := bot.session.Open(); err != nil {
		return errors.Wrap(err, "could not open session")
	}

	return nil
}

func (bot *DiscordBot) Dispose() error {
	defer bot.session.Close()

	if bot.RemoveCommands {
		bot.log.Info("Removing commands...")

		for _, cmd := range bot.commands {
			err := bot.session.ApplicationCommandDelete(bot.session.State.User.ID, bot.GuildID, cmd.ID)
			if err != nil {
				return errors.Wrapf(err, "Cannot delete '%v'", cmd.Name)
			}
		}
	}

	return nil
}
