package predictor

import (
	"context"
	"log"

	"github.com/carlmjohnson/requests"
	"github.com/karatekaneen/stockybot"
)

type Predictor struct {
	URL string
}

type PredictionRequest struct {
	StockData          []stockybot.PricePoint `json:"stock_data,omitempty"`
	OmxData            []stockybot.PricePoint `json:"omx_data,omitempty"`
	TradesThisYear     int                    `json:"trades_this_year,omitempty"`
	DaysSinceLastTrade int                    `json:"days_since_last_trade,omitempty"`
}

type PredictionResponse struct {
	Prediction float64
}

func (p *Predictor) SignalRank(ctx context.Context, req PredictionRequest) (*PredictionResponse, error) {
	// var resp PredictionResponse
	var s string
	// Yes, we can be lazy and use a dependency here
	err := requests.
		URL("/predict").
		Host(p.URL).
		BodyJSON(&req).ToString(&s).
		// ToJSON(&resp).
		Fetch(ctx)

	log.Println(s)

	return nil, err
}
