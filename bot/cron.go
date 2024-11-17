package bot

import (
	"context"
	_ "embed"
	"fmt"
	"slices"
	"strings"
	"text/template"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"github.com/karatekaneen/stockybot"
)

//go:embed templates/dailyreport.tmpl
var reportTemplate string

type cronController struct {
	watchRepo subscriptionRepository
	log       *zap.SugaredLogger
	ranker    *rankController
	session   *discordgo.Session
	schedule  string
	cfg       Config
}

func newCronController(
	cfg Config,
	log *zap.SugaredLogger,
	session *discordgo.Session,
	ranker *rankController,
	watchRepo subscriptionRepository,
) *cronController {
	cc := &cronController{
		cfg:       cfg,
		log:       log.Named("cron").With("schedule", cfg.Schedule),
		session:   session,
		ranker:    ranker,
		watchRepo: watchRepo,
	}

	go cc.dailyReport()

	c := cron.New()
	c.AddFunc(cfg.Schedule, func() { cc.dailyReport() })

	return cc
}

func (cc *cronController) dailyReport() {
	cc.log.Info("Starting daily report")

	ctx, cancel := context.WithTimeout(context.Background(), cc.cfg.DefaultTimeout)
	defer cancel()

	summary, err := cc.generateDailyReport(ctx)
	if err != nil {
		cc.log.Errorw("Generate daily report", "err", err)
		cc.session.ChannelMessageSend(
			cc.cfg.ChannelID,
			fmt.Sprintf("generate daily report: %s", err),
		)
	}

	strSummary, err := printSummary(reportTemplate, *summary)
	if err != nil {
		cc.log.Errorw("Compile report message", "err", err)
		cc.session.ChannelMessageSend(
			cc.cfg.ChannelID,
			fmt.Sprintf("generate daily report string: %s", err),
		)
	}

	if _, err := cc.session.ChannelMessageSend(cc.cfg.ChannelID, strSummary); err != nil {
		cc.log.Errorw("Send message", "err", err)
	}
}

func (rc *cronController) generateDailyReport(ctx context.Context) (*dailySummary, error) {
	indexPrices, err := rc.ranker.getIndexPrices(ctx)
	if err != nil {
		return nil, fmt.Errorf("get index prices: %w", err)
	}

	pending, err := rc.ranker.getGroupedPendingSigs(ctx)
	if err != nil {
		return nil, fmt.Errorf("get pending buys: %w", err)
	}

	allSignals := []int64{}
	for _, sig := range pending.sells {
		allSignals = append(allSignals, sig.Stock.ID)
	}

	for _, sig := range pending.buys {
		allSignals = append(allSignals, sig.Stock.ID)
	}

	watches, err := rc.watchRepo.GetWatchersBySecurities(ctx, allSignals)
	if err != nil {
		return nil, fmt.Errorf("get watchers by securities: %w", err)
	}

	preds, err := rc.ranker.makePredictions(ctx, pending.buys, indexPrices)
	if err != nil {
		err = wrapErr(err, "make predictions: %w")
		return nil, err
	}

	summary := relevantSignals(pending, watches, preds)

	slices.SortFunc(summary.Buys, asDescendingScore)

	return &summary, nil
}

type watchSignal struct {
	Watchers []string
	prediction
}

func newWatchSignal(sig stockybot.Signal, watchers []string, score float64) watchSignal {
	return watchSignal{
		prediction: prediction{
			Signal: sig,
			score:  score,
		},
		Watchers: watchers,
	}
}

type dailySummary struct {
	Sells []watchSignal
	Buys  []watchSignal
}

func relevantSignals(
	pending *groupedSignals,
	watchMap map[int64][]string,
	preds []prediction,
) dailySummary {
	out := dailySummary{
		Sells: []watchSignal{},
		Buys:  []watchSignal{},
	}

	for _, buy := range pending.buys {
		watchers := watchMap[buy.Stock.ID]
		score := 0.0

		predIdx := slices.IndexFunc(
			preds,
			func(p prediction) bool { return p.Signal.Stock.ID == buy.Stock.ID },
		)
		if predIdx >= 0 {
			score = preds[predIdx].score
		}

		if len(watchers) > 0 || score >= 0.3 {
			out.Buys = append(out.Buys, newWatchSignal(buy, watchers, score))
		}
	}

	for _, sell := range pending.sells {
		if watchers := watchMap[sell.Stock.ID]; len(watchers) > 0 {
			out.Sells = append(out.Sells, newWatchSignal(sell, watchMap[sell.Stock.ID], 0))
		}
	}

	return out
}

func printSummary(rawReportTmpl string, summary dailySummary) (string, error) {
	// TODO: This should probably not be parsed everytime
	reportTmpl, err := template.New("foo").Parse(rawReportTmpl)
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}

	strB := &strings.Builder{}

	err = reportTmpl.Execute(strB, summary)
	return strB.String(), err
}
