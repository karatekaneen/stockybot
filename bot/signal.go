package bot

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/karatekaneen/stockybot"
	"github.com/karatekaneen/stockybot/predictor"
	"github.com/pkg/errors"
)

type signalGroup struct {
	buys  []stockybot.Signal
	sells []stockybot.Signal
}

// TODO: Maybe split up?
type dataRepository interface {
	PendingSignals(ctx context.Context) ([]stockybot.Signal, error)
	Signals(ctx context.Context, stockId int) ([]stockybot.Signal, error)
	PriceData(ctx context.Context, id int) (*stockybot.PriceDocument, error)
	Security(ctx context.Context, id int) (*stockybot.Security, error)
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
	ctx, cancel := context.WithTimeout(context.Background(), bot.defaultTimeout)
	defer cancel()

	// Access options in the order provided by the user.
	optionMap := mapOptions(i.ApplicationCommandData().Options)

	var allLists bool
	if opt, ok := optionMap["all-lists"]; ok {
		allLists = opt.BoolValue()
	}

	bot.log.Info("Starting to fetch signals")

	signals, err := bot.signalRepository.PendingSignals(ctx)
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

func (bot *DiscordBot) rankBuySignals(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx, cancel := context.WithTimeout(context.Background(), bot.defaultTimeout)
	defer cancel()

	bot.log.Info("Starting to fetch signals")

	// BTS: 5503
	stockId := 1051357
	// stockId := 19002

	// Omxs30: 19002
	omxId := 19002
	// omxId := 1051357

	indexPriceDoc, err := bot.signalRepository.PriceData(ctx, omxId)
	if err != nil {
		bot.log.Errorln(errors.Wrap(err, "Stock price fetch:"))
		interactionResponse(s, i, "An error occured")
		return
	}

	stockPriceDoc, err := bot.signalRepository.PriceData(ctx, stockId)
	if err != nil {
		bot.log.Errorln(errors.Wrap(err, "Stock price fetch:"))
		interactionResponse(s, i, "An error occured")
		return
	}

	signals, err := bot.signalRepository.Signals(ctx, stockId)
	if err != nil {
		bot.log.Errorln(errors.Wrap(err, "Signal price fetch:"))
		interactionResponse(s, i, "An error occured")
		return
	}

	stockPrices, err := stockybot.LastN(stockPriceDoc.PriceData, 201)
	if err != nil {
		bot.log.Errorln(errors.Wrap(err, "Not enough data:"))
		interactionResponse(s, i, "An error occured")
		return
	}
	omxPrices, err := stockybot.LastN(indexPriceDoc.PriceData, 201)
	if err != nil {
		bot.log.Errorln(errors.Wrap(err, "Not enough data:"))
		interactionResponse(s, i, "An error occured")
		return
	}

	p := predictor.Predictor{URL: "stockyml-dudb2aklkq-lz.a.run.app"}

	prediction, err := p.SignalRank(ctx, predictor.PredictionRequest{
		StockData:          stockPrices,
		OmxData:            omxPrices,
		TradesThisYear:     numberOfExitsSince(signals, time.Now().Add(-time.Hour*24*365)),
		DaysSinceLastTrade: daysSinceLast(signals, time.Now()),
	})
	if err != nil {
		bot.log.Errorln(errors.Wrap(err, "Prediction request:"))
		interactionResponse(s, i, "An error occured")
		return
	}

	log.Printf("%+v", prediction)

	bot.log.Info("Starting sending response")
	if err := interactionResponse(s, i, "hell to the yes"); err != nil {
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
			// fmt.Sprintf(	"Yolo  <@%s> %v", i.Member.User.ID, allLists,),
		},
	})
}
