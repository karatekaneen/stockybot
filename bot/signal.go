package bot

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"

	"github.com/karatekaneen/stockybot"
)

type signalGroup struct {
	buys  []stockybot.Signal
	sells []stockybot.Signal
}

// TODO: Maybe split up?
type dataRepository interface {
	PendingSignals(ctx context.Context) ([]stockybot.Signal, error)
	Signals(ctx context.Context, stockId int64) ([]stockybot.Signal, error)
	PriceData(ctx context.Context, id int64) (*stockybot.PriceDocument, error)
	Security(ctx context.Context, id int64) (*stockybot.Security, error)
}

func mapOptions(
	options []*discordgo.ApplicationCommandInteractionDataOption,
) map[string]*discordgo.ApplicationCommandInteractionDataOption {
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	return optionMap
}

func (bot *DiscordBot) pendingSignals(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx, cancel := context.WithTimeout(context.Background(), bot.cfg.DefaultTimeout)
	defer cancel()

	// Access options in the order provided by the user.
	optionMap := mapOptions(i.ApplicationCommandData().Options)

	var allLists bool
	if opt, ok := optionMap["all-lists"]; ok {
		allLists = opt.BoolValue()
	}

	bot.log.Info("Starting to fetch signals")

	signals, err := bot.dataRepository.PendingSignals(ctx)
	if err != nil {
		bot.log.Errorln(errors.Wrap(err, "Pending signal fetch:"))
		interactionResponse(s, i, "An error occured")
		return
	}
	bot.log.Infow("Fetched signals", "signals", len(signals))

	// Group signals by list and type
	groupedSignals := groupSignals(signals)

	// Create strings for each list and type
	content := signalsByListAndType(groupedSignals, bot.defaultStockLists, allLists)

	bot.log.Info("Starting sending response")
	if err := interactionResponse(s, i, strings.Join(content, "\n")); err != nil {
		bot.log.Error(err)
		return
	}

	bot.log.Info("Sent response")
}

func daysSinceLast(signals []stockybot.Signal, now time.Time) int {
	if len(signals) < 1 {
		return 0
	}

	return int(now.Sub(signals[len(signals)-1].Date) / time.Hour / 24)
}

func numberOfExitsSince(signals []stockybot.Signal, limit time.Time) int {
	numSignals := 0
	for _, sig := range signals {
		if sig.Type == "exit" && sig.Date.After(limit) {
			numSignals++
		}
	}

	return numSignals
}

func signalsByListAndType(
	groupedSignals map[string]signalGroup,
	defaultStockLists map[string]struct{},
	allLists bool,
) []string {
	content := []string{}
	format := "\n**%s - %s**\n\t%s"

	for list, signals := range groupedSignals {
		if _, defaultList := defaultStockLists[list]; !defaultList && !allLists {
			continue // k thx bye
		}

		if len(signals.buys) > 0 {
			content = append(content, fmt.Sprintf(
				format,
				list,
				"BUY",
				strings.Join(stockNames(signals.buys), ", "),
			))
		}
		if len(signals.sells) > 0 {
			content = append(content, fmt.Sprintf(
				format,
				list,
				"SELL",
				strings.Join(stockNames(signals.sells), ", "),
			))
		}
	}

	// Add message if no signals found
	if len(content) < 1 {
		content = append(content, "No signals found for the selected lists")
	}

	return content
}

func groupSignals(signals []stockybot.Signal) map[string]signalGroup {
	groupedSignals := map[string]signalGroup{}
	for _, signal := range signals {
		signalList, ok := groupedSignals[signal.Stock.List]
		if !ok {
			signalList = signalGroup{}
		}

		if signal.Action == "buy" {
			signalList.buys = append(signalList.buys, signal)
		} else {
			signalList.sells = append(signalList.sells, signal)
		}

		groupedSignals[signal.Stock.List] = signalList
	}

	return groupedSignals
}

// Extracts stock names from signals.
// TODO: Sort alphabetically
func stockNames(signals []stockybot.Signal) []string {
	stockNames := []string{}

	for _, signal := range signals {
		stockNames = append(stockNames, signal.Stock.Name)
	}

	return stockNames
}

func followUpResponse(
	s *discordgo.Session,
	i *discordgo.Interaction,
	content string,
) error {
	_, err := s.FollowupMessageCreate(i, true, &discordgo.WebhookParams{Content: content})
	return err
	// return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
	// 	// Ignore type for now, they will be discussed in "responses"
	// 	Type: discordgo.InteractionResponseChannelMessageWithSource,
	// 	Data: &discordgo.InteractionResponseData{
	// 		Content: content,
	// 	},
	// })
}

func interactionResponse(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	content string,
) error {
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		// Ignore type for now, they will be discussed in "responses"
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
}
