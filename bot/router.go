package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/karatekaneen/stockybot/predictor"
	"go.uber.org/zap"
)

// FIXME:
// - watch ska ha autocomplete på add och remove

type interactionHandler func(s *discordgo.Session, i *discordgo.InteractionCreate) error

type action struct {
	command  *discordgo.ApplicationCommand
	handlers map[discordgo.InteractionType]interactionHandler
	name     string
}

func (a *action) Command() *discordgo.ApplicationCommand {
	return a.command
}

type router struct {
	logger    *zap.SugaredLogger
	ranker    *rankController
	signaller *signalController
	commands  []action
}

func newRouter(
	config Config,
	log *zap.SugaredLogger,
	repo dataRepository,
	pred *predictor.Predictor,
) *router {
	return &router{
		logger: log.Named("router"),
		signaller: &signalController{
			log:            log,
			cfg:            config,
			dataRepository: repo,
			defaultStockLists: map[string]struct{}{
				"Small Cap Stockholm":  {},
				"Mid Cap Stockholm":    {},
				"Large Cap Stockholm":  {},
				"Large Cap Copenhagen": {},
			},
		},
		ranker: &rankController{
			log:            log,
			cfg:            config,
			dataRepository: repo,
			predictor:      pred,
		},
	}
}

func (r *router) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	logger := r.logger.With(
		"command", i.ApplicationCommandData().Name,
		"type", i.Type,
	)

	for _, foo := range r.commands {
		if foo.name != i.ApplicationCommandData().Name {
			continue
		}

		handler, ok := foo.handlers[i.Type]
		if !ok {
			break
		}

		if err := handler(s, i); err != nil {
			logger.Error(err)
		}
	}

	logger.Error("No handler found")
}

func (r *router) Commands() []*discordgo.ApplicationCommand {
	cmds := []*discordgo.ApplicationCommand{}

	for _, cmd := range r.actions() {
		cmds = append(cmds, cmd.Command())
	}

	return cmds
}

func (r *router) actions() []action {
	return []action{
		{
			name: "pending",
			handlers: map[discordgo.InteractionType]interactionHandler{
				discordgo.InteractionApplicationCommand: r.signaller.pendingSignals,
			},
			command: &discordgo.ApplicationCommand{
				Name:        "pending",
				Description: "List pending signals",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionBoolean,
						Name:        "all-lists",
						Description: "Show signals from all lists. Only listing Swedish Large, Mid and Small cap if false",
					},
				},
			},
		},
		{
			name: "rankbuys",
			handlers: map[discordgo.InteractionType]interactionHandler{
				discordgo.InteractionApplicationCommand: r.ranker.rankBuySignals,
			},
			command: &discordgo.ApplicationCommand{
				Name:        "rankbuys",
				Description: "List pending buy signals by signal predicted ranking",
			},
		},
		{
			// FIXME: Add handlers here
			name: "watch",
			command: &discordgo.ApplicationCommand{
				Name:        "watch",
				Description: "Add or remove subscriptions of stocks",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "list",
						Description: "List your subscriptions",
						Type:        discordgo.ApplicationCommandOptionSubCommand,
					},
					{
						Name:        "add",
						Description: "Add subscription of a stock",
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Options: []*discordgo.ApplicationCommandOption{
							{
								Type:        discordgo.ApplicationCommandOptionString,
								Name:        "ticker",
								Description: "Ticker of the stock you want to subscribe to",
								Required:    true,
							},
						},
					},
					{
						Name:        "remove",
						Description: "Remove subscription of a stock",
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Options: []*discordgo.ApplicationCommandOption{
							{
								Type:        discordgo.ApplicationCommandOptionString,
								Name:        "ticker",
								Description: "Ticker of the stock you want to subscribe to",
								Required:    true,
							},
						},
					},
				},
			},
		},
	}
}
