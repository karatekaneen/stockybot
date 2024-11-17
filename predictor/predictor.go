package predictor

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/carlmjohnson/requests"
	"go.uber.org/zap"

	"github.com/karatekaneen/stockybot"
)

type Predictor struct {
	Config
	log *zap.SugaredLogger
}
type Config struct {
	PredictURL string `help:"URL to prediction model service" required:"" env:"PREDICTION_URL"`
	ScaleURL   string `help:"URL to data scaling service"     required:"" env:"SCALING_URL"`
}

func New(cfg Config, log *zap.SugaredLogger) *Predictor {
	return &Predictor{cfg, log.Named("predictor")}
}

type PredictionRequest struct {
	StockData          []stockybot.PricePoint `json:"stock_data,omitempty"`
	OmxData            []stockybot.PricePoint `json:"omx_data,omitempty"`
	TradesThisYear     int                    `json:"trades_this_year"`
	DaysSinceLastTrade int                    `json:"days_since_last_trade"`
}

type PredictionResponse struct {
	Predictions [][]float64 `json:"predictions"`
}

func (p *Predictor) scale(ctx context.Context, req PredictionRequest) ([][]float64, error) {
	resp := make([][]float64, 1)

	rawResp := &bytes.Buffer{}

	err := requests.
		URL(fmt.Sprint(p.ScaleURL, "/scale")).
		BodyJSON(&req).
		ToBytesBuffer(rawResp).
		Fetch(ctx)
	if err != nil {
		if requests.HasStatusErr(err, http.StatusTooManyRequests) {
			time.Sleep(time.Second)
			return p.scale(ctx, req)
		}
		return nil, fmt.Errorf("scale data: %s: %w", rawResp.String(), err)
	}

	if err := json.Unmarshal(rawResp.Bytes(), &resp); err != nil {
		if strings.Contains(rawResp.String(), "NaN") {
			return nil, errors.Join(stockybot.ErrInvalidData, err)
		}
		return nil, fmt.Errorf("json unmarshal: %q, %w", rawResp.String(), err)
	}

	return resp, nil
}

func (p *Predictor) predict(
	ctx context.Context,
	req [][]float64,
) (*PredictionResponse, error) {
	var resp PredictionResponse

	err := requests.
		URL(fmt.Sprint(p.PredictURL, "/v1/models/stocky/versions/1:predict")).
		BodyJSON(map[string][][]float64{"instances": req}).
		ToJSON(&resp).
		Fetch(ctx)
	if err != nil {
		if requests.HasStatusErr(err, http.StatusTooManyRequests) {
			time.Sleep(time.Second)
			return p.predict(ctx, req)
		}
		return nil, fmt.Errorf("prediction failed: %w", err)
	}

	return &resp, nil
}

func (p *Predictor) SignalRank(
	ctx context.Context,
	req PredictionRequest,
) (float64, error) {
	// var resp PredictionResponse
	scaled, err := p.scale(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("scale data: %w", err)
	}

	pred, err := p.predict(ctx, scaled)
	if err != nil {
		return 0, fmt.Errorf("predcition: %w", err)
	}

	if len(pred.Predictions) == 0 || len(pred.Predictions[0]) == 0 {
		return 0, fmt.Errorf("unexpected shape")
	}

	return pred.Predictions[0][0], nil
}
