package bot

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	dscd "github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/karatekaneen/stockybot"
	"github.com/karatekaneen/stockybot/predictor"
)

type prediction struct {
	Signal stockybot.Signal
	score  float64
}

func (p prediction) Score() float64 { return p.score * 100 }

type rankController struct {
	log            *zap.SugaredLogger
	cfg            Config
	dataRepository dataRepository
	predictor      *predictor.Predictor
}

func (rc *rankController) rankBuySignals(s *dscd.Session, i *dscd.InteractionCreate) error {
	ctx, cancel := context.WithTimeout(context.Background(), rc.cfg.DefaultTimeout)
	defer cancel()

	rc.log.Info("Starting to fetch signals")
	interactionResponse(s, i, "Making predictions for buy signals, ignoring under 30%...")

	// TODO: In the future you could make this and the request below concurrent
	indexPrices, err := rc.getIndexPrices(ctx)
	if err != nil {
		err = wrapErr(err, "get index prices: %w")
		return err
	}

	pending, err := rc.getGroupedPendingSigs(ctx)
	if err != nil {
		err = wrapErr(err, "get pending buys: %w")
		interactionErr(s, i, err)
		return err
	}

	preds, err := rc.makePredictions(ctx, pending.buys, indexPrices)
	if err != nil {
		err = wrapErr(err, "make predictions: %w")
		interactionErr(s, i, err)
		return err
	}

	slices.SortFunc(preds, asDescendingScore)

	summaries := make([]string, 0, len(preds))

	for _, p := range preds {
		if p.score < 0.3 {
			continue
		}

		summaries = append(
			summaries,
			fmt.Sprintf("- %.1f%% - %s", p.Score(), p.Signal.Stock.Name),
		)
	}

	if err := followUpResponse(s, i.Interaction, strings.Join(summaries, "\n")); err != nil {
		return fmt.Errorf("send followup response: %w", err)
	}

	return nil
}

func (rc *rankController) getIndexPrices(ctx context.Context) ([]stockybot.PricePoint, error) {
	indexPriceDoc, err := rc.dataRepository.PriceData(ctx, rc.cfg.IndexID)
	if err != nil {
		return nil, fmt.Errorf("get index price data: %w", err)
	}

	indexPrices, err := stockybot.LastN(indexPriceDoc.PriceData, 201)
	if err != nil {
		return nil, fmt.Errorf("get last 201 items of index data: %w", err)
	}

	return indexPrices, nil
}

type groupedSignals struct {
	buys, sells []stockybot.Signal
}

func (rc *rankController) getGroupedPendingSigs(ctx context.Context) (*groupedSignals, error) {
	pendingSignals, err := rc.dataRepository.PendingSignals(ctx)
	if err != nil {
		return nil, wrapErr(err, "get pending signals: %w")
	}

	pending := &groupedSignals{buys: []stockybot.Signal{}, sells: []stockybot.Signal{}}

	for _, sig := range pendingSignals {
		if sig.Action == "buy" {
			pending.buys = append(pending.buys, sig)
		} else {
			pending.sells = append(pending.sells, sig)
		}
	}

	return pending, nil
}

func (rc *rankController) makePredictions(
	ctx context.Context,
	pending []stockybot.Signal,
	indexPrices []stockybot.PricePoint,
) ([]prediction, error) {
	preds := make([]prediction, 0, len(pending))
	mut := &sync.Mutex{}

	rc.log.Infof("Fetching predictions for %d pending signals", len(pending))

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(8)

	for _, pendingSig := range pending {
		pendingSig := pendingSig // TODO: remove in 1.22

		g.Go(func() error {
			req, err := rc.createPredictionRequest(ctx, indexPrices, pendingSig.Stock.ID)
			if err != nil {
				rc.log.Errorw(
					"an error occured while making predictions. Ignoring.",
					"id", pendingSig.Stock.ID,
					"error", err,
				)
				return nil
			}

			log := rc.log.With(
				"id", pendingSig.Stock.ID,
				"omx", len(req.OmxData),
				"stock", len(req.StockData),
				"thisYear", req.TradesThisYear,
				"days", req.DaysSinceLastTrade,
			)

			log.Info("making prediction")

			predictionScore, err := rc.predictor.SignalRank(ctx, *req)
			if err != nil {
				log.Errorw("an error occured while making predictions. Ignoring.", "error", err)
			} else {
				log.Info("Got prediction OK")
			}

			mut.Lock()
			preds = append(preds, prediction{Signal: pendingSig, score: predictionScore})
			mut.Unlock()

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("make predictions: %w", err)
	}

	return preds, nil
}

func (rc *rankController) createPredictionRequest(
	ctx context.Context,
	indexPrices []stockybot.PricePoint,
	stockId int64,
) (*predictor.PredictionRequest, error) {
	stockPriceDoc, err := rc.dataRepository.PriceData(ctx, stockId)
	if err != nil {
		return nil, fmt.Errorf("fetch prices for security: %d: %w", stockId, err)
	}

	signals, err := rc.dataRepository.Signals(ctx, stockId)
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

type scorer interface {
	Score() float64
}

func asDescendingScore[T scorer](a, b T) int {
	switch {
	case a.Score() < b.Score():
		return 1
	case a.Score() > b.Score():
		return -1
	default:
		return 0
	}
}
