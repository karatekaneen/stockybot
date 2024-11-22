package bot

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"

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
	StrategyState(ctx context.Context, id int64) (*stockybot.StrategyState, error)
}

type signalController struct {
	log               *zap.SugaredLogger
	dataRepository    dataRepository
	defaultStockLists map[string]struct{}
	cfg               Config
}

func (sc *signalController) pendingSignals(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), sc.cfg.DefaultTimeout)
	defer cancel()

	// Access options in the order provided by the user.
	optionMap := mapOptions(i.ApplicationCommandData().Options)

	var allLists bool
	if opt, ok := optionMap["all-lists"]; ok {
		allLists = opt.BoolValue()
	}

	sc.log.Info("Starting to fetch signals")

	signals, err := sc.dataRepository.PendingSignals(ctx)
	if err != nil {
		interactionResponse(s, i, "An error occured")
		return fmt.Errorf("get pending signals: %w", err)
	}

	sc.log.Infow("Fetched signals", "signals", len(signals))

	// Group signals by list and type
	groupedSignals := groupSignals(signals)

	// Create strings for each list and type
	content := signalsByListAndType(groupedSignals, sc.defaultStockLists, allLists)

	sc.log.Info("Starting sending response")
	if err := interactionResponse(s, i, strings.Join(content, "\n")); err != nil {
		return fmt.Errorf("send response: %w", err)
	}

	return nil
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
