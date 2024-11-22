package firestore

import (
	"context"
	"fmt"
	"strconv"

	"github.com/karatekaneen/stockybot"
)

type rawStrategyState struct {
	Bias         string     `firestore:"bias,omitempty"`
	Regime       string     `firestore:"regime,omitempty"`
	LastSignal   rawSignal2 `firestore:"lastSignal,omitempty"`
	HighPrice    float64    `firestore:"highPrice,omitempty"`
	LowPrice     float64    `firestore:"lowPrice,omitempty"`
	TriggerPrice float64    `firestore:"triggerPrice,omitempty"`
}

func (f *FireDB) StrategyState(ctx context.Context, id int64) (*stockybot.StrategyState, error) {
	doc, err := f.client.
		Collection("context").
		Doc(strconv.FormatInt(id, 10)).
		Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetch strategy context document: %w", err)
	} else if !doc.Exists() {
		return nil, ErrNotExist
	}

	var raw rawStrategyState
	if err := doc.DataTo(&raw); err != nil {
		return nil, fmt.Errorf("convert context document: %w", err)
	}

	sig, err := raw.LastSignal.toSignal()
	if err != nil {
		return nil, fmt.Errorf("convert raw signal to proper signal: %w", err)
	}

	out := stockybot.StrategyState{
		Bias:         raw.Bias,
		Regime:       raw.Regime,
		LastSignal:   sig,
		HighPrice:    raw.HighPrice,
		LowPrice:     raw.LowPrice,
		TriggerPrice: raw.TriggerPrice,
	}

	return &out, nil
}
