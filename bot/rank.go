package bot

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/sync/errgroup"

	"github.com/karatekaneen/stockybot"
	"github.com/karatekaneen/stockybot/predictor"
)

type prediction struct {
	signal stockybot.Signal
	score  float64
}

func (bot *DiscordBot) rankBuySignals(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx, cancel := context.WithTimeout(context.Background(), bot.cfg.DefaultTimeout)
	defer cancel()

	bot.log.Info("Starting to fetch signals")
	interactionResponse(s, i, "Making predictions for buy signals, ignoring under 30%...")

	// TODO: In the future you could make this and the request below concurrent

	indexPrices, err := bot.getIndexPrices(ctx)
	if err != nil {
		bot.interactionErr(s, i, wrapErr(err, "get index prices: %w"))
		return
	}

	pendingBuys, err := bot.getPendingBuys(ctx)
	if err != nil {
		bot.interactionErr(s, i, wrapErr(err, "get pending buys: %w"))
		return
	}

	preds, err := bot.makePredictions(ctx, pendingBuys, indexPrices)
	if err != nil {
		bot.interactionErr(s, i, wrapErr(err, "make predictions: %w"))
		return
	}

	slices.SortFunc(preds, asDescending)

	summaries := make([]string, 0, len(preds))

	for _, p := range preds {
		if p.score < 0.3 {
			continue
		}

		summaries = append(
			summaries,
			fmt.Sprintf("- %.1f%% - %s", p.score*100, p.signal.Stock.Name),
		)
	}

	if err := followUpResponse(s, i.Interaction, strings.Join(summaries, "\n")); err != nil {
		bot.log.Error(err)
		return
	}
	bot.log.Info("Sent response")
}

func (bot *DiscordBot) getIndexPrices(ctx context.Context) ([]stockybot.PricePoint, error) {
	indexPriceDoc, err := bot.dataRepository.PriceData(ctx, bot.cfg.IndexID)
	if err != nil {
		return nil, fmt.Errorf("get index price data: %w", err)
	}

	indexPrices, err := stockybot.LastN(indexPriceDoc.PriceData, 201)
	if err != nil {
		return nil, fmt.Errorf("get last 201 items of index data: %w", err)
	}

	return indexPrices, nil
}

func (bot *DiscordBot) interactionErr(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	err error,
) error {
	bot.log.Error(err)
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: "An error occured"},
	})
}

func wrapErr(err error, wrapper string) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf(wrapper, err)
}

func (bot *DiscordBot) getPendingBuys(ctx context.Context) ([]stockybot.Signal, error) {
	pendingSignals, err := bot.dataRepository.PendingSignals(ctx)
	if err != nil {
		return nil, wrapErr(err, "get pending signals: %w")
	}

	pendingBuys := []stockybot.Signal{}

	for _, sig := range pendingSignals {
		if sig.Action == "buy" {
			pendingBuys = append(pendingBuys, sig)
		}
	}

	return pendingBuys, nil
}

func (bot *DiscordBot) makePredictions(
	ctx context.Context,
	pending []stockybot.Signal,
	indexPrices []stockybot.PricePoint,
) ([]prediction, error) {
	preds := make([]prediction, 0, len(pending))
	mut := &sync.Mutex{}

	bot.log.Infof("Fetching predictions for %d pending signals", len(pending))

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(8)

	for _, pendingSig := range pending {
		pendingSig := pendingSig // TODO: remove in 1.22

		g.Go(func() error {
			req, err := bot.createPredictionRequest(ctx, indexPrices, pendingSig.Stock.ID)
			if err != nil {
				bot.log.Errorw(
					"an error occured while making predictions. Ignoring.",
					"id", pendingSig.Stock.ID,
					"error", err,
				)
				return nil
			}

			log := bot.log.With(
				"id", pendingSig.Stock.ID,
				"omx", len(req.OmxData),
				"stock", len(req.StockData),
				"thisYear", req.TradesThisYear,
				"days", req.DaysSinceLastTrade,
			)

			log.Info("making prediction")

			predictionScore, err := bot.predictor.SignalRank(ctx, *req)
			if err != nil {
				log.Errorw("an error occured while making predictions. Ignoring.", "error", err)
			} else {
				log.Info("Got prediction OK")
			}

			mut.Lock()
			preds = append(preds, prediction{signal: pendingSig, score: predictionScore})
			mut.Unlock()

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("make predictions: %w", err)
	}

	return preds, nil
}

func (bot *DiscordBot) createPredictionRequest(
	ctx context.Context,
	indexPrices []stockybot.PricePoint,
	stockId int64,
) (*predictor.PredictionRequest, error) {
	stockPriceDoc, err := bot.dataRepository.PriceData(ctx, stockId)
	if err != nil {
		return nil, fmt.Errorf("fetch prices for security: %d: %w", stockId, err)
	}

	signals, err := bot.dataRepository.Signals(ctx, stockId)
	if err != nil && errors.Is(err, stockybot.ErrNotFound) {
		return nil, fmt.Errorf("fetch signals for security: %d: %w", stockId, err)
	}

	stockPrices, err := stockybot.LastN(stockPriceDoc.PriceData, 201)
	if err != nil {
		return nil, fmt.Errorf("not enough data for %d: %w", stockId, err)
	}

	return &predictor.PredictionRequest{
		StockData:          stockPrices,
		OmxData:            indexPrices,
		TradesThisYear:     numberOfExitsSince(signals, time.Now().Add(-time.Hour*24*365)),
		DaysSinceLastTrade: daysSinceLast(signals, time.Now()),
	}, nil
}

func asDescending(a, b prediction) int {
	switch {
	case a.score < b.score:
		return 1
	case a.score > b.score:
		return -1
	default:
		return 0
	}
}
