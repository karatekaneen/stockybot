package bot

import (
	"fmt"

	dscd "github.com/bwmarrin/discordgo"
	"go.uber.org/zap"

	"github.com/karatekaneen/stockybot/predictor"
)

// FIXME:
// - watch ska ha autocomplete på add och remove

type interactionHandler func(s *dscd.Session, i *dscd.InteractionCreate) error

type action struct {
	command *dscd.ApplicationCommand
	// The interaction type can be thought of similar to a
	// HTTP method such as POST, GET, etc. It will define
	// what kind of interaction that will happen. The handler
	// is then responsible of doing it.
	handlers map[dscd.InteractionType]interactionHandler
	name     string
}

func (a *action) Command() *dscd.ApplicationCommand {
	return a.command
}

type router struct {
	logger    *zap.SugaredLogger
	ranker    *rankController
	watcher   *watchController
	signaller *signalController
}

func newRouter(
	config Config,
	log *zap.SugaredLogger,
	repo dataRepository,
	pred *predictor.Predictor,
	subRepo subscriptionRepository,
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
		watcher: &watchController{
			log:       log,
			watchRepo: subRepo,
		},
		ranker: &rankController{
			log:            log,
			cfg:            config,
			dataRepository: repo,
			predictor:      pred,
		},
	}
}

func (r *router) Handle(s *dscd.Session, i *dscd.InteractionCreate) {
	logger := r.logger.With(
		"command", i.ApplicationCommandData().Name,
		"type", i.Type,
	)

	for _, foo := range r.actions() {
		if foo.name != i.ApplicationCommandData().Name {
			continue
		}

		handler, ok := foo.handlers[i.Type]
		if !ok {
			break
		}

		if err := handler(s, i); err != nil {
			logger.Error(err)

			// Send the error to the user
			errContent := fmt.Sprintf(
				"An error occurred when performing command %s: %v",
				i.ApplicationCommandData().Name,
				err,
			)
			if err := interactionResponse(s, i, errContent); err != nil {
				logger.Error(err)
			}
		}
		return
	}

	logger.Error("No handler found")

	if err := interactionResponse(s, i, "No event handler found"); err != nil {
		logger.Error(err)
	}
}

func (r *router) Commands() []*dscd.ApplicationCommand {
	cmds := []*dscd.ApplicationCommand{}

	for _, cmd := range r.actions() {
		cmds = append(cmds, cmd.Command())
	}

	return cmds
}

func (r *router) actions() []action {
	return []action{
		{
			name: "pending",
			handlers: map[dscd.InteractionType]interactionHandler{
				dscd.InteractionApplicationCommand: r.signaller.pendingSignals,
			},
			command: &dscd.ApplicationCommand{
				Name:        "pending",
				Description: "List pending signals",
				Options: []*dscd.ApplicationCommandOption{
					{
						Type:        dscd.ApplicationCommandOptionBoolean,
						Name:        "all-lists",
						Description: "Show signals from all lists. Only listing Swedish Large, Mid and Small cap if false",
					},
				},
			},
		},

		{
			name: "rankbuys",
			handlers: map[dscd.InteractionType]interactionHandler{
				dscd.InteractionApplicationCommand: r.ranker.rankBuySignals,
			},
			command: &dscd.ApplicationCommand{
				Name:        "rankbuys",
				Description: "List pending buy signals by signal predicted ranking",
			},
		},

		{
			name: "watch-list",
			handlers: map[dscd.InteractionType]interactionHandler{
				dscd.InteractionApplicationCommand: r.watcher.List,
			},
			command: &dscd.ApplicationCommand{
				Name:        "watch-list",
				Description: "List your subscriptions",
			},
		},

		{
			name: "watch-add",
			handlers: map[dscd.InteractionType]interactionHandler{
				dscd.InteractionApplicationCommand:             r.watcher.AddCommit,
				dscd.InteractionApplicationCommandAutocomplete: r.watcher.AddAutocomplete,
			},
			command: &dscd.ApplicationCommand{
				Name:        "watch-add",
				Description: "Add a new subscription of a stock",
				Options: []*dscd.ApplicationCommandOption{
					{
						Type:         dscd.ApplicationCommandOptionString,
						Name:         "ticker",
						Description:  "Ticker of the stock you want to subscribe to",
						Required:     true,
						Autocomplete: true,
					},
				},
			},
		},

		{
			name: "watch-remove",
			handlers: map[dscd.InteractionType]interactionHandler{
				dscd.InteractionApplicationCommand:             r.watcher.RemoveCommit,
				dscd.InteractionApplicationCommandAutocomplete: r.watcher.RemoveAutocomplete,
			},
			command: &dscd.ApplicationCommand{
				Name:        "watch-remove",
				Description: "remove a subscription of stocks",
				Options: []*dscd.ApplicationCommandOption{
					{
						Type:         dscd.ApplicationCommandOptionString,
						Name:         "ticker",
						Description:  "Ticker of the stock you want to subscribe to",
						Required:     true,
						Autocomplete: true,
					},
				},
			},
		},
	}
}
